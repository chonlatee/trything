package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/chonlatee/payment/internal/lib/typemapper"
	"golang.org/x/tools/imports"
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
}

func (g *gen) run() error {
	g.el = make([]entityInfo, 0)
	err := filepath.WalkDir("./internal/dbgen", g.walkFile)
	if err != nil {
		return err
	}

	g.generateEntity()
	g.generateModelToEntity()

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
			srcType := typemapper.AstExprToString(v.Type)
			e.fields[i] = entityField{
				name:    v.Names[0].Name,
				srcType: srcType,
				dstType: typemapper.PgxTypeMapper(srcType),
			}
		}

		g.el = append(g.el, e)

		return true
	})

	return nil
}

func (g *gen) generateEntity() {
	t, err := os.ReadFile("./internal/tpl/entity.tpl")
	if err != nil {
		log.Fatal(err)
	}

	fnc := template.FuncMap{
		"getEntityList": func() []entityInfo {
			return g.el
		},
		"getEntityName": func(e entityInfo) string {
			return e.name
		},
		"getEntityFields": func(e entityInfo) []entityField {
			return e.fields
		},
		"getEntityFieldName": func(e entityField) string {
			return e.name
		},
		"getEntityFieldType": func(e entityField) string {
			return e.dstType
		},
	}

	tpl := template.Must(template.New("entity").Funcs(fnc).Parse(string(t)))

	var buf bytes.Buffer

	err = tpl.Execute(&buf, g.el)
	if err != nil {
		log.Fatal(err)
	}

	r, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	out, err := imports.Process("entity.go", r, &imports.Options{
		Comments:  true,
		TabWidth:  4,
		TabIndent: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./internal/entity/entity.go", out, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("generate entity success.")
}

func (g *gen) generateModelToEntity() {

}
