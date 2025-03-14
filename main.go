package main

import (
	_ "embed"
	"errors"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

//go:embed enum.go.template
var templateEnum string

//go:embed enum_test.go.template
var templateEnumTest string

//go:embed enum_string.go.template
var templateString string

//go:embed enum_string_test.go.template
var templateStringTest string

func main() {
	var (
		typeName     string
		enableString bool
		fileName     = os.Getenv("GOFILE")
		lineNum      = os.Getenv("GOLINE")
		packageName  = os.Getenv("GOPACKAGE")
	)
	flag.StringVar(&typeName, "type", "", "type to be generated for")
	flag.BoolVar(&enableString, "string", false, "generate String() method")
	flag.Parse()

	if err := process(packageName, fileName, typeName, lineNum, enableString); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

func process(packageName, fileName, typeName, lineNum string, enableString bool) error {
	if typeName == "" || fileName == "" || packageName == "" || lineNum == "" {
		return errors.New("missing parameters")
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

	expectedLine, _ := strconv.Atoi(lineNum)
	expectedLine += 1

	var specs [][2]string

	ast.Inspect(f, func(astNode ast.Node) bool {
		node, ok := astNode.(*ast.GenDecl)
		if !ok || node == nil || (node.Tok != token.CONST && node.Tok != token.VAR) {
			return true
		}

		if position := fset.Position(node.Pos()); position.Line != expectedLine {
			return false
		}

		for _, astSpec := range node.Specs {
			spec, ok := astSpec.(*ast.ValueSpec)
			if !ok || spec == nil {
				continue
			}

			if len(spec.Names) != 1 {
				break
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
		return errors.New(fileName + ": unable to find values for enum type: " + typeName)
	}

	templateCode := templateEnum
	templateTest := templateEnumTest
	if enableString {
		templateCode += "\n" + templateString
		templateTest += "\n" + templateStringTest
	}

	replacer := strings.NewReplacer(
		"{{.Type}}", typeName,
		"{{.Package}}", packageName,
		"{{.Values}}", mapJoin(specs, func(_ int, v [2]string) string { return v[0] }, ", "),
		"{{.Tags}}", mapJoin(specs, func(_ int, v [2]string) string { return `"` + v[1] + `"` }, ","),
		"{{.TagsNaked}}", mapJoin(specs, func(_ int, v [2]string) string { return v[1] }, " "),
		"{{.seq_bytes}}", mapJoin(specs, func(_ int, v [2]string) string { return `[]byte("` + v[1] + `")` }, ", "),
		"{{.seq_string}}", mapJoin(specs, func(_ int, v [2]string) string { return `"` + v[1] + `"` }, ", "),
		"{{.string_to_value_switch}}", mapJoin(specs, func(_ int, v [2]string) string { return `case "` + v[1] + "\":\n *s = " + v[0] }, "\n"),
		"{{.value_to_append_bytes_switch}}", mapJoin(specs, func(i int, v [2]string) string {
			return `case ` + v[0] + ":\n return append(b, seq_bytes_" + typeName + "[" + strconv.Itoa(i) + `]...), nil`
		}, "\n"),
		"{{.value_to_string_switch}}", mapJoin(specs, func(i int, v [2]string) string {
			return `case ` + v[0] + ":\n return seq_string_" + typeName + "[" + strconv.Itoa(i) + `]`
		}, "\n"),
	)

	bastPath := filepath.Join(filepath.Dir(fileName), camelCaseToSnakeCase(typeName))

	return errors.Join(
		writeCode([]byte(replacer.Replace(templateCode)), bastPath+"_enum_encoding.go"),
		writeCode([]byte(replacer.Replace(templateTest)), bastPath+"_enum_encoding_test.go"),
	)
}

func mapJoin(vs [][2]string, f func(idx int, val [2]string) string, sep string) string {
	var b strings.Builder
	for i, v := range vs {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(f(i, v))
	}
	return b.String()
}

func writeCode(code []byte, outFilePath string) error {
	formattedCode, err := format.Source(code)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilePath, formattedCode, 0644)
}

func camelCaseToSnakeCase(s string) string {
	var b strings.Builder
	var isPrevUpper bool
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && !isPrevUpper {
				b.WriteRune('_')
			}
			isPrevUpper = true
			b.WriteRune(unicode.ToLower(r))
		} else {
			isPrevUpper = false
			b.WriteRune(r)
		}
	}
	return b.String()
}
