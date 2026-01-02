package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var fset = token.NewFileSet()

type modelInfo struct {
	name  string
	field []modelField
}

type modelField struct {
	name  string
	ftype ast.Expr
}

type modelDataFields struct {
	name      string
	fieldType string
}

type modelData struct {
	modelName string
	fields    []modelDataFields
}

func (s *st) walkFile(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if filepath.Ext(path) == ".go" {
		file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}

			for _, spec := range gen.Specs {
				ts := spec.(*ast.TypeSpec)
				if st, ok := ts.Type.(*ast.StructType); ok {
					if ts.Name.Name != "Queries" {
						m := modelInfo{
							name:  ts.Name.Name,
							field: make([]modelField, 0),
						}
						for _, field := range st.Fields.List {
							m.field = append(m.field, modelField{
								name:  field.Names[0].Name,
								ftype: field.Type,
							})
						}
						s.modelInfo = append(s.modelInfo, m)
					}
				}
			}
		}

	}

	return nil
}

type st struct {
	dir       string
	modelInfo []modelInfo
}

func main() {
	s := st{
		dir:       "./dbgen",
		modelInfo: []modelInfo{},
	}

	if err := s.run(); err != nil {
		log.Fatal(err)
	}

}

func (s *st) run() error {
	err := filepath.WalkDir(s.dir, s.walkFile)
	if err != nil {
		return err
	}

	s.generateModel()
	return err
}

func (s *st) generateModel() {

	tmpl := `package entity`

	tmpl += `type {{getModelName .}}Model struct {
		{{ range getModelFields .}}
		{{getFieldName .}}	{{getFieldType .}}
		{{end}}
	}`

	for _, v := range s.modelInfo {
		m := modelData{
			modelName: v.name,
			fields:    make([]modelDataFields, 0),
		}
		for _, f := range v.field {
			m.fields = append(m.fields, modelDataFields{
				name:      f.name,
				fieldType: mapFieldType(exprToString(f.ftype)),
			})
		}

	}
	funcMaps := template.FuncMap{
		"getModelName":   func(m modelData) string { return m.modelName },
		"getModelFields": func(m modelData) []modelDataFields { return m.fields },
		"getFieldName":   func(m modelDataFields) string { return m.name },
		"getFieldType":   func(m modelDataFields) string { return m.fieldType },
	}
	tpl := template.Must(template.New("tmpl").Funcs(funcMaps).Parse(tmpl))
	f, err := os.ReadFile("./tmpl/entity.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.Execute(f, m)
	if err != nil {
		log.Printf("tpl excute err: %v\n", err)
	}
}

func mapFieldType(f string) string {
	if v, ok := mapType[f]; ok {
		return v
	}

	return f
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		if pkg, ok := t.X.(*ast.Ident); ok {
			return fmt.Sprintf("%s.%s", pkg.Name, t.Sel.Name)
		}
		return t.Sel.Name
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", exprToString(t.Key), exprToString(t.Value))
	case *ast.InterfaceType:
		return "any"
	default:
		return fmt.Sprintf("%T", expr)
	}
}
