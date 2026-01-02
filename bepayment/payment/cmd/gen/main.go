package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
)

var fset = token.NewFileSet()

type gen struct {
	el []entityInfo
}

type entityInfo struct {
	name   string
	fields []entityField
}

type entityField struct {
	name    string
	srcType string
	dstType string
}

func main() {
	g := &gen{}

	if err := g.run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("entity list: %+v\n", g.el)
}

func (g *gen) run() error {
	g.el = make([]entityInfo, 0)
	err := filepath.WalkDir("./internal/dbgen", g.walkFile)
	if err != nil {
		return err
	}

	return nil
}

func (g *gen) walkFile(path string, dir fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if filepath.Ext(path) != ".go" {
		return nil
	}

	file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	ast.Inspect(file, func(n ast.Node) bool {

		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if ts.Name.Name == "Queries" {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		e := entityInfo{
			name:   ts.Name.Name,
			fields: make([]entityField, len(st.Fields.List)),
		}

		for i, v := range st.Fields.List {
			srcType := astExprToString(v.Type)
			e.fields[i] = entityField{
				name:    v.Names[0].Name,
				srcType: srcType,
				dstType: pgxTypeMapper(srcType),
			}
		}

		g.el = append(g.el, e)

		return true
	})

	return nil
}
