package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed enum.go.template
var templateCode string

//go:embed enum_short.go.template
var templateShortCode string

//go:embed enum_test.go.template
var templateTest string

func main() {
	var (
		typeName    string
		mode        string
		fileName    = os.Getenv("GOFILE")
		packageName = os.Getenv("GOPACKAGE")
	)
	flag.StringVar(&typeName, "type", "", "type to be generated for")
	flag.StringVar(&mode, "mode", "auto", "what kind of strategy used (short, long, auto)")
	flag.Parse()

	if err := process(typeName, fileName, packageName, mode); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

func process(typeName string, fileName string, packageName string, mode string) error {
	if typeName == "" || fileName == "" || packageName == "" {
		return errors.New("type, file and package name must be provided")
	}

	inputCode, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, inputCode, parser.ParseComments)
	if err != nil {
		return err
	}

	var specs [][2]string

	ast.Inspect(f, func(astNode ast.Node) bool {
		node, ok := astNode.(*ast.GenDecl)
		if !ok || (node.Tok != token.CONST && node.Tok != token.VAR) {
			return true
		}

		var typeFound *ast.Ident

		for _, astSpec := range node.Specs {
			spec, ok := astSpec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			if len(spec.Names) != 1 {
				break
			}

			if typeFound == nil {
				if spec.Type != nil {
					// type stated explicitly; e.g.:
					//   `Green Color = iota + 1 // json:"green"`
					typeFound, _ = spec.Type.(*ast.Ident)
				} else if valueFound, ok := spec.Values[0].(*ast.CompositeLit); ok {
					// type not stated, but is composite (expected struct); e.g.:
					//   `Green Color = Color{1} // json:"green"`
					typeFound, ok = valueFound.Type.(*ast.Ident)
				} else {
					// type not stated, and is either a literal or an expression,
					// impossible to infer the enum type name; e.g.:
					//   `Green = iota + 1 // json:"green"`
					position := fset.Position(spec.Pos())
					os.Stderr.WriteString(fmt.Sprintf(
						"%s: (warning) unable to infer enum type; " +
						"Add explicit type or use struct type.\n",
						position))
				}
				if typeFound == nil || typeFound.Name != typeName {
					break
				}
			}

			tag, ok := "", false
			for _, field := range strings.Fields(spec.Comment.Text()) {
				if strings.HasPrefix(field, "json:") {
					tag, ok = field[len("json:\""):len(field)-1], true
					break
				}
			}
			if ok {
				specs = append(specs, [2]string{spec.Names[0].Name, tag})
			}
		}

		return false
	})

	if len(specs) == 0 {
		return errors.New(
			fmt.Sprintf("%s: Unable to find values for enum type %q\n", fileName, typeName))
	}

	code := templateCode

	if mode == "auto" {
		mode = "short"
		if len(specs) >= 10 {
			mode = "long"
		}
	}

	if mode == "short" {
		code = templateShortCode
		code = strings.ReplaceAll(code, "{{.json_to_value}}", strings.Join(mp(specs, func(_ int, v [2]string) string { return `case "` + v[1] + "\":\n *s = " + v[0] }), "\n"))
		code = strings.ReplaceAll(code, "{{.val_to_json}}", strings.Join(mp(specs, func(i int, v [2]string) string {
			return `case ` + v[0] + ":\n return json_bytes_" + typeName + "[" + strconv.Itoa(i) + `], nil`
		}), "\n"))
		code = strings.ReplaceAll(code, "{{.json_bytes}}", strings.Join(mp(specs, func(_ int, v [2]string) string { return `[]byte("` + v[1] + `")` }), ", "))
	} else {
		code = strings.ReplaceAll(code, "{{.val_to_json}}", strings.Join(mp(specs, func(_ int, v [2]string) string { return v[0] + `: []byte("` + v[1] + `"),` }), "\n"))
		code = strings.ReplaceAll(code, "{{.json_to_value}}", strings.Join(mp(specs, func(_ int, v [2]string) string { return `"` + v[1] + `": ` + v[0] + `,` }), "\n"))
	}

	code = strings.ReplaceAll(code, "{{.Type}}", typeName)
	code = strings.ReplaceAll(code, "{{.Package}}", packageName)

	test := templateTest
	test = strings.ReplaceAll(test, "{{.Type}}", typeName)
	test = strings.ReplaceAll(test, "{{.Package}}", packageName)
	test = strings.ReplaceAll(test, "{{.Values}}", strings.Join(mp(specs, func(_ int, v [2]string) string { return v[0] }), ", "))
	test = strings.ReplaceAll(test, "{{.Tags}}", strings.Join(mp(specs, func(_ int, v [2]string) string { return `"` + v[1] + `"` }), ","))

	return errors.Join(
		writeCode([]byte(code), filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding.go")),
		writeCode([]byte(test), filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding_test.go")),
	)
}

func mp[T any, M any](a []T, f func(int, T) M) (l []M) {
	for i, e := range a {
		l = append(l, f(i, e))
	}
	return l
}

func writeCode(code []byte, outFilePath string) error {
	formattedCode, err := format.Source(code)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilePath, formattedCode, 0644)
}
