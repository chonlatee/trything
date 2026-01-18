package typemapper

import (
	"go/ast"
	"testing"
)

func Test_pgxTypeMapper(t *testing.T) {

	testCases := []struct {
		name string
		src  string
		dst  string
	}{
		{
			name: "string to string",
			src:  "string",
			dst:  "string",
		},
		{
			name: "int to int",
			src:  "int",
			dst:  "int",
		},
		{
			name: "pgtype.UUID to string",
			src:  "pgtype.UUID",
			dst:  "string",
		},
		{
			name: "[]pgtype.UUID to []string",
			src:  "[]pgtype.UUID",
			dst:  "[]string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PgxTypeMapper(tc.src)
			if got != tc.dst {
				t.Errorf("got %s want %s", got, tc.dst)
			}
		})
	}
}

func Test_astExprToString(t *testing.T) {
	testCases := []struct {
		name string
		expr ast.Expr
		want string
	}{
		{
			name: "Ident int to int",
			expr: ast.NewIdent("int"),
			want: "int",
		},
		{
			name: "Ident string to string",
			expr: ast.NewIdent("string"),
			want: "string",
		},
		{
			name: "StarExpr string to *string",
			expr: &ast.StarExpr{
				X: ast.NewIdent("string"),
			},
			want: "*string",
		},
		{
			name: "StarExpr int to *int",
			expr: &ast.StarExpr{
				X: ast.NewIdent("int"),
			},
			want: "*int",
		},
		{
			name: "ArrayType string to []string",
			expr: &ast.ArrayType{
				Elt: ast.NewIdent("string"),
			},
			want: "[]string",
		},
		{
			name: "ArrayType *string to []*string",
			expr: &ast.ArrayType{
				Elt: &ast.StarExpr{
					X: ast.NewIdent("string"),
				},
			},
			want: "[]*string",
		},
		{
			name: "MapType [string]string to map[string]string",
			expr: &ast.MapType{
				Key:   ast.NewIdent("string"),
				Value: ast.NewIdent("string"),
			},
			want: "map[string]string",
		},
		{
			name: "MapType [string]*string to map[string]*string",
			expr: &ast.MapType{
				Key: ast.NewIdent("string"),
				Value: &ast.StarExpr{
					X: ast.NewIdent("string"),
				},
			},
			want: "map[string]*string",
		},
		{
			name: "Selector pgx.UUID to pgx.UUID",
			expr: &ast.SelectorExpr{
				X:   ast.NewIdent("pgx"),
				Sel: ast.NewIdent("UUID"),
			},
			want: "pgx.UUID",
		},
		{
			name: "Selector *pgx.UUID to *pgx.UUID",
			expr: &ast.SelectorExpr{
				X: &ast.StarExpr{
					X: ast.NewIdent("pgx"),
				},
				Sel: ast.NewIdent("UUID"),
			},
			want: "*pgx.UUID",
		},
		{
			name: "ArrayType Selector []pgx.UUID to []pgx.UUID",
			expr: &ast.ArrayType{
				Elt: &ast.SelectorExpr{
					X:   ast.NewIdent("pgx"),
					Sel: ast.NewIdent("UUID"),
				},
			},
			want: "[]pgx.UUID",
		},
		{
			name: "ArrayType StarExpr Selector []*pgx.UUID to []*pgx.UUID",
			expr: &ast.ArrayType{
				Elt: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("pgx"),
						Sel: ast.NewIdent("UUID"),
					},
				},
			},
			want: "[]*pgx.UUID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := AstExprToString(tc.expr)
			if tc.want != got {
				t.Errorf("got %s want %s", got, tc.want)
			}
		})
	}
}
