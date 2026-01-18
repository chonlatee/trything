package typemapper

import (
	"fmt"
	"go/ast"
	"strings"
)

var pgxtype = map[string]string{
	"pgtype.UUID":        "string",
	"pgtype.Numeric":     "float64",
	"pgtype.Timestamptz": "time.Time",
	"pgtype.Bool":        "bool",
	"pgtype.Text":        "string",
}

func PgxTypeMapper(src string) string {
	if v, ok := pgxtype[src]; ok {
		return v
	}

	// []pgtype.xxx
	if strings.HasPrefix(src, "[]") {
		return "[]" + PgxTypeMapper(src[2:])
	}

	return src
}

func AstExprToString(ex ast.Expr) string {
	switch t := ex.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + AstExprToString(t.X)
	case *ast.ArrayType:
		return "[]" + AstExprToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", AstExprToString(t.Key), AstExprToString(t.Value))
	case *ast.SelectorExpr:
		return AstExprToString(t.X) + "." + t.Sel.Name
	default:
		return fmt.Sprintf("%s", ex)
	}
}
