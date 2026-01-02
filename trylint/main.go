package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"unicode"
)

type visitor struct {
	fset *token.FileSet
}

func main() {

	v := &visitor{fset: token.NewFileSet()}

	for _, filePath := range os.Args[1:] {
		if filePath == "--" {
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}

		ast.Walk(v, f)
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}

	name := funcDecl.Name.String()

	n := splitCamelCase(name)

	if len(n) == 2 {
		if _, ok := EntityWords[n[1]]; ok {
			if len(funcDecl.Type.Results.List) == 1 {
				if e, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident); ok {
					if e.Name == "string" {
						fmt.Printf("%s: function name '%s' should has identify in name like `%s(Name|ID|Title)` when return single value.\n", v.fset.Position(node.Pos()), funcDecl.Name.Name, funcDecl.Name.Name)
					}
				}
			}
		}
	}

	// resultList := funcDecl.Type.Results.List

	// log.Printf("func name: %s\n", funcDecl.Name.String())

	// var buf bytes.Buffer
	// printer.Fprint(&buf, v.fset, node)

	// fmt.Printf("%s | %#v\n", buf.String(), node)

	return v
}

func splitCamelCase(s string) []string {
	var words []string

	current := ""

	for _, r := range s {
		if unicode.IsUpper(r) && current != "" {
			words = append(words, current)
			current = ""
		}

		current += string(r)
	}

	if current != "" {
		words = append(words, current)
	}

	return words
}
