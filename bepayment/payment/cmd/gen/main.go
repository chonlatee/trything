package main

import (
	"bytes"
	"fmt"
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
	g.generateEntityToModel()

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
		"getModelName": func(e entityInfo) string {
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
	t, err := os.ReadFile("./internal/tpl/entityToModel.tpl")
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
		"getModelName": func(e entityInfo) string {
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
		"getEntityFieldSrcType": func(e entityField) string {
			return e.srcType
		},
		"getEntityFieldValue": func(e entityField) string {
			if e.dstType == "[]string" && e.srcType == "[]pgtype.UUID" {
				return fmt.Sprintf("converter.StringsToPgTypeUUIDs(input.%s)", e.name)
			}

			if e.dstType == "[]float64" && e.srcType == "[]pgtype.Numeric" {
				return fmt.Sprintf("converter.Float64sToPgTypeNumerics(input.%s)", e.name)
			}

			if e.dstType == "[]time.Time" && e.srcType == "[]pgtype.Timestamptz" {
				return fmt.Sprintf("converter.TimesToPgTypeTimes(input.%s)", e.name)
			}

			if e.dstType == "[]string" && e.srcType == "[]string" {
				return fmt.Sprintf("input.%s", e.name)
			}

			if e.dstType == "string" && e.srcType == "pgtype.UUID" {
				return fmt.Sprintf("converter.StringToPgtypeUUID(input.%s)", e.name)
			}

			if e.dstType == "float64" && e.srcType == "pgtype.Numeric" {
				return fmt.Sprintf("converter.Float64ToPgtypeNumeric(input.%s)", e.name)
			}

			if e.dstType == "string" && e.srcType == "pgtype.Text" {
				return fmt.Sprintf("converter.StringToPgtypeText(input.%s)", e.name)
			}

			if e.dstType == "time.Time" && e.srcType == "pgtype.Timestamptz" {
				return fmt.Sprintf("converter.TimeToPgtypeTimestamptz(input.%s)", e.name)
			}

			return fmt.Sprintf("input.%s", e.name)
		},

		"needConvert": func(e entityField) bool {
			return e.dstType != e.srcType
		},
	}

	tpl := template.Must(template.New("entityToModel").Funcs(fnc).Parse(string(t)))

	var buf bytes.Buffer

	err = tpl.Execute(&buf, g.el)
	if err != nil {
		log.Fatal(err)
	}

	r, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	out, err := imports.Process("entitiToModel.go", r, &imports.Options{
		Comments:  true,
		TabWidth:  4,
		TabIndent: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./internal/entity/entityToModel.go", out, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("generate entity to model success.")
}

func (g *gen) generateEntityToModel() {
	t, err := os.ReadFile("./internal/tpl/modelToEntity.tpl")
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
		"getModelName": func(e entityInfo) string {
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
		"getEntityFieldSrcType": func(e entityField) string {
			return e.srcType
		},
		"getModelFieldValue": func(e entityField) string {
			if e.dstType == "[]string" && e.srcType == "[]pgtype.UUID" {
				return fmt.Sprintf("converter.PgtypeUUIDsToStrings(input.%s)", e.name)
			}

			if e.dstType == "[]float64" && e.srcType == "[]pgtype.Numeric" {
				return fmt.Sprintf("converter.PgtypeNumericsToFloat64s(input.%s)", e.name)
			}

			if e.dstType == "[]time.Time" && e.srcType == "[]pgtype.Timestamptz" {
				return fmt.Sprintf("converter.PgtypeTimestamptzsToTimes(input.%s)", e.name)
			}

			if e.dstType == "[]string" && e.srcType == "[]string" {
				return fmt.Sprintf("input.%s", e.name)
			}

			if e.dstType == "string" && e.srcType == "pgtype.UUID" {
				return fmt.Sprintf("converter.PgtypeUUIDToString(input.%s)", e.name)
			}

			if e.dstType == "float64" && e.srcType == "pgtype.Numeric" {
				return fmt.Sprintf("converter.PgtypeNumericToFloat64(input.%s)", e.name)
			}

			if e.dstType == "string" && e.srcType == "pgtype.Text" {
				return fmt.Sprintf("converter.PgtypeTextToString(input.%s)", e.name)
			}

			if e.dstType == "time.Time" && e.srcType == "pgtype.Timestamptz" {
				return fmt.Sprintf("converter.PgtypeTimestamptzToTime(input.%s)", e.name)
			}

			return fmt.Sprintf("input.%s", e.name)
		},

		"needConvert": func(e entityField) bool {
			return e.dstType != e.srcType
		},
	}

	tpl := template.Must(template.New("modelToEntity").Funcs(fnc).Parse(string(t)))

	var buf bytes.Buffer

	err = tpl.Execute(&buf, g.el)
	if err != nil {
		log.Fatal(err)
	}

	r, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	out, err := imports.Process("modelToEntity.go", r, &imports.Options{
		Comments:  true,
		TabWidth:  4,
		TabIndent: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./internal/entity/modelToEntity.go", out, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("generate model to entity success.")

}
