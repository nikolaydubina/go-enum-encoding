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

var (
	typeName    string
	fileName    = os.Getenv("GOFILE")
	packageName = os.Getenv("GOPACKAGE")
)

func main() {
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
	code = bytes.ReplaceAll(code, []byte("{{.val_to_json}}"), []byte(strings.Join(valToJson(specs), "\n")))
	code = bytes.ReplaceAll(code, []byte("{{.json_to_value}}"), []byte(strings.Join(jsonToVal(specs), "\n")))

	test := templateTest
	test = bytes.ReplaceAll(test, []byte("{{.Type}}"), []byte(typeName))
	test = bytes.ReplaceAll(test, []byte("{{.Package}}"), []byte(packageName))
	test = bytes.ReplaceAll(test, []byte("{{.Values}}"), []byte(strings.Join(vals(specs), ", ")))
	test = bytes.ReplaceAll(test, []byte("{{.Tags}}"), []byte(strings.Join(tags(specs), ",")))

	return errors.Join(
		writeCode(code, filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding.go")),
		writeCode(test, filepath.Join(filepath.Dir(fileName), strings.ToLower(typeName)+"_enum_encoding_test.go")),
	)
}

func valToJson(specs map[string]string) (res []string) {
	for val, jsonTag := range specs {
		res = append(res, fmt.Sprintf("%s: \"%s\",", val, jsonTag))
	}
	sort.StringSlice(res).Sort()
	return res
}

func jsonToVal(specs map[string]string) (res []string) {
	for val, jsonTag := range specs {
		res = append(res, fmt.Sprintf("\"%s\": %s,", jsonTag, val))
	}
	sort.StringSlice(res).Sort()
	return res
}

func vals(specs map[string]string) (res []string) {
	for val := range specs {
		res = append(res, val)
	}
	sort.StringSlice(res).Sort()
	return res
}

func tags(specs map[string]string) (res []string) {
	for _, tag := range specs {
		res = append(res, `"`+tag+`"`)
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
