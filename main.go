package main

import (
	_ "embed"
	"errors"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed enum.go.template
var templateCode string

//go:embed enum_test.go.template
var templateTest string

func main() {
	var (
		typeName    string
		fileName    = os.Getenv("GOFILE")
		packageName = os.Getenv("GOPACKAGE")
	)
	flag.StringVar(&typeName, "type", "", "type to be generated for")
	flag.Parse()

	if err := process(typeName, fileName, packageName); err != nil {
		log.Fatalf("cannot process: %s", err)
	}
}

func process(typeName string, fileName string, packageName string) error {
	if typeName == "" || fileName == "" || packageName == "" {
		return errors.New("type, file and package name must be provided")
	}

	inputCode, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	f, err := parser.ParseFile(token.NewFileSet(), fileName, inputCode, parser.ParseComments)
	if err != nil {
		return err
	}

	specs := make(map[string]string)

	ast.Inspect(f, func(node ast.Node) bool {
		spec, ok := node.(*ast.ValueSpec)
		if !ok {
			return true
		}

		if len(spec.Names) != 1 {
			return false
		}
		// TODO: check that type matches

		var jsonTag string
		for _, field := range strings.Fields(spec.Comment.Text()) {
			if strings.HasPrefix(field, "json:") {
				jsonTag = field[len("json:\"") : len(field)-1]
				break
			}
		}
		if jsonTag == "" {
			return false
		}

		specs[spec.Names[0].Name] = jsonTag
		return false
	})

	if tag, ok := specs["Undefined"+typeName]; !ok {
		return errors.New(`missing enum symbol("Undefined` + typeName + `") for type(` + typeName + `)`)
	} else if tag != "-" {
		return errors.New(`json tag for Undefined value must be "-" but got "` + tag + `"`)
	}
	delete(specs, "Undefined"+typeName)

	code := templateCode
	code = strings.ReplaceAll(code, "{{.Type}}", typeName)
	code = strings.ReplaceAll(code, "{{.Package}}", packageName)
	code = strings.ReplaceAll(code, "{{.val_to_json}}", strings.Join(mp(specs, func(k, v string) string { return k + `: "` + v + `",` }), "\n"))
	code = strings.ReplaceAll(code, "{{.json_to_value}}", strings.Join(mp(specs, func(k, v string) string { return `"` + v + `": ` + k + `,` }), "\n"))

	test := templateTest
	test = strings.ReplaceAll(test, "{{.Type}}", typeName)
	test = strings.ReplaceAll(test, "{{.Package}}", packageName)
	test = strings.ReplaceAll(test, "{{.Values}}", strings.Join(mp(specs, func(k, _ string) string { return k }), ", "))
	test = strings.ReplaceAll(test, "{{.Tags}}", strings.Join(mp(specs, func(_, v string) string { return `"` + v + `"` }), ","))

	return errors.Join(
		writeCode([]byte(code), filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding.go")),
		writeCode([]byte(test), filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding_test.go")),
	)
}

func mp(m map[string]string, f func(k, v string) string) (res []string) {
	for k, v := range m {
		res = append(res, f(k, v))
	}
	sort.StringSlice(res).Sort()
	return res
}

func writeCode(code []byte, outFilePath string) error {
	formattedCode, err := format.Source(code)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilePath, formattedCode, 0644)
}
