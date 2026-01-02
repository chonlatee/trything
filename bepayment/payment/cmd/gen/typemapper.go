package main

import (
	"fmt"
	"go/ast"
	"strings"
)

var pgxtype = map[string]string{
	"pgtype.UUID":        "string",
	"pgtype.Numeric":     "float64",
	"pgtype.Timestamptz": "time.Timestamptz",
	"pgtype.Bool":        "bool",
}

func pgxTypeMapper(src string) string {
	if v, ok := pgxtype[src]; ok {
		return v
	}

	// []pgtype.xxx
	if strings.HasPrefix(src, "[]") {
		return "[]" + pgxTypeMapper(src[2:])
	}

	return src
}

func astExprToString(ex ast.Expr) string {
	switch t := ex.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + astExprToString(t.X)
	case *ast.ArrayType:
		return "[]" + astExprToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", astExprToString(t.Key), astExprToString(t.Value))
	case *ast.SelectorExpr:
		return astExprToString(t.X) + "." + t.Sel.Name
	default:
		return fmt.Sprintf("%s", ex)
	}
}
