package yourlint

import (
	"go/ast"
	"go/token"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("yourlinter", New)
}

type YourSettings struct{}

type Yourlinter struct {
	settings YourSettings
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[YourSettings](settings)
	if err != nil {
		return nil, err
	}

	return &Yourlinter{settings: s}, nil
}

func (y *Yourlinter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "yourlinter",
			Doc:  "your linter",
			Run:  y.run,
		},
	}, nil
}

func (y *Yourlinter) GetLoadMode() string {
	return register.LoadModeSyntax
}

func (y *Yourlinter) run(pass *analysis.Pass) (any, error) {

	hasString := make(map[string]bool)

	structPos := make(map[string]token.Pos)

	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {

			switch x := n.(type) {
			case *ast.TypeSpec:
				if _, ok := x.Type.(*ast.StructType); ok {
					structPos[x.Name.Name] = x.Pos()

					if _, exists := hasString[x.Name.Name]; !exists {
						hasString[x.Name.Name] = false
					}
				}
			case *ast.FuncDecl:
				if x.Recv != nil && len(x.Recv.List) > 0 {
					typeName := getReceiverTypeName(x.Recv.List[0].Type)
					if x.Name.Name == "String" && typeName != "" {
						hasString[typeName] = true
					}
				}

			}

			return true
		})
	}

	for name, found := range hasString {
		if !found {
			pass.Reportf(structPos[name], "struct %s is missing a String() method", name)
		}
	}

	return nil, nil
}

func getReceiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}

	return ""
}
