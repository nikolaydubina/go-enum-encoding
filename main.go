package main

import (
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
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
var templateCode []byte

//go:embed enum_test.go.template
var templateTest []byte

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
		return fmt.Errorf("type, file and package name must be provided")
	}

	inputCode, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot read file(%s): %s", fileName, err)
	}

	f, err := parser.ParseFile(token.NewFileSet(), fileName, inputCode, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("cannot parse file(%s): %s", fileName, err)
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
		return fmt.Errorf(`missing enum symbol("%s") for type(%s)`, "Undefined"+typeName, typeName)
	} else if tag != "-" {
		return fmt.Errorf(`json tag for Undefined value must be "-" but got "%s"`, tag)
	}
	delete(specs, "Undefined"+typeName)

	code := templateCode
	code = bytes.ReplaceAll(code, []byte("{{.Type}}"), []byte(typeName))
	code = bytes.ReplaceAll(code, []byte("{{.Package}}"), []byte(packageName))
	code = bytes.ReplaceAll(code, []byte("{{.val_to_json}}"), []byte(strings.Join(mp(specs, func(k, v string) string { return fmt.Sprintf("%s: \"%s\",", k, v) }), "\n")))
	code = bytes.ReplaceAll(code, []byte("{{.json_to_value}}"), []byte(strings.Join(mp(specs, func(k, v string) string { return fmt.Sprintf("\"%s\": %s,", v, k) }), "\n")))

	test := templateTest
	test = bytes.ReplaceAll(test, []byte("{{.Type}}"), []byte(typeName))
	test = bytes.ReplaceAll(test, []byte("{{.Package}}"), []byte(packageName))
	test = bytes.ReplaceAll(test, []byte("{{.Values}}"), []byte(strings.Join(mp(specs, func(k, _ string) string { return k }), ", ")))
	test = bytes.ReplaceAll(test, []byte("{{.Tags}}"), []byte(strings.Join(mp(specs, func(_, v string) string { return `"` + v + `"` }), ",")))

	return errors.Join(
		writeCode(code, filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding.go")),
		writeCode(test, filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding_test.go")),
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
		return fmt.Errorf("cannot format code: %s", err)
	}
	return os.WriteFile(outFilePath, formattedCode, 0644)
}
