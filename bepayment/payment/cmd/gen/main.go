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
	return filepath.WalkDir("./internal/dbgen", g.walkFile)

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
			srcType := g.exprToString(v.Type)
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

func (g *gen) exprToString(ex ast.Expr) string {
	switch t := ex.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + g.exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + g.exprToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", g.exprToString(t.Key), g.exprToString(t.Value))
	case *ast.SelectorExpr:
		if pkg, ok := t.X.(*ast.Ident); ok {
			return fmt.Sprintf("%s.%s", pkg.Name, t.Sel.Name)
		}
		return t.Sel.Name
	default:
		return fmt.Sprintf("%s", ex)
	}
}
