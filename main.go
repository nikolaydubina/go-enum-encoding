package main

import (
	"bytes"
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
		if !ok || (node.Tok != token.CONST && node.Tok != token.VAR) {
			return true
		}

		position := fset.Position(node.Pos())
		if position.Line != expectedLine {
			return false
		}

		for _, astSpec := range node.Specs {
			spec, ok := astSpec.(*ast.ValueSpec)
			if !ok {
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

	r := (&replacer{vals: make(map[string]string), specs: specs}).
		With("{{.Type}}", typeName).
		With("{{.Package}}", packageName).
		WithMap("{{.Values}}", func(_ int, v [2]string) string { return v[0] }, ", ").
		WithMap("{{.Tags}}", func(_ int, v [2]string) string { return `"` + v[1] + `"` }, ",").
		WithMap("{{.TagsNaked}}", func(_ int, v [2]string) string { return v[1] }, " ").
		WithMap("{{.seq_bytes}}", func(_ int, v [2]string) string { return `[]byte("` + v[1] + `")` }, ", ").
		WithMap("{{.seq_string}}", func(_ int, v [2]string) string { return `"` + v[1] + `"` }, ", ")

	templateCode := templateEnum
	templateTest := templateEnumTest
	if enableString {
		templateCode += "\n" + templateString
		templateTest += "\n" + templateStringTest
	}

	r = r.
		WithMap("{{.string_to_value_switch}}", func(_ int, v [2]string) string { return `case "` + v[1] + "\":\n *s = " + v[0] }, "\n").
		WithMap("{{.value_to_bytes_switch}}", func(i int, v [2]string) string {
			return `case ` + v[0] + ":\n return seq_bytes_" + typeName + "[" + strconv.Itoa(i) + `], nil`
		}, "\n").
		WithMap("{{.value_to_string_switch}}", func(i int, v [2]string) string {
			return `case ` + v[0] + ":\n return seq_string_" + typeName + "[" + strconv.Itoa(i) + `]`
		}, "\n")

	bastPath := filepath.Join(filepath.Dir(fileName), camelCaseToSnakeCase(strings.ToLower(typeName)))

	return errors.Join(
		writeCode(r.Apply([]byte(templateCode)), bastPath+"_enum_encoding.go"),
		writeCode(r.Apply([]byte(templateTest)), bastPath+"_enum_encoding_test.go"),
	)
}

type replacer struct {
	vals  map[string]string
	specs [][2]string
}

func (r *replacer) With(k, v string) *replacer {
	r.vals[k] = v
	return r
}

func (r *replacer) WithMap(k string, f func(idx int, val [2]string) string, sep string) *replacer {
	strings.NewReplacer()

	var b strings.Builder
	for i, v := range r.specs {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(f(i, v))
	}
	r.vals[k] = b.String()
	return r
}

func (r *replacer) Apply(s []byte) []byte {
	for k, v := range r.vals {
		s = bytes.ReplaceAll(s, []byte(k), []byte(v))
	}
	return []byte(s)
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
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && !unicode.IsUpper(rune(s[i-1])) {
				b.WriteRune('_')
			}
			r = unicode.ToLower(r)
		}
		b.WriteRune(r)
	}
	return b.String()
}
