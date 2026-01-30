package mylinter

import (
	"go/ast"
	"strings"
	"unicode"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("mylinter", New)
}

type MySettings struct{}

type MyLinter struct {
	settings MySettings
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[MySettings](settings)
	if err != nil {
		return nil, err
	}

	return &MyLinter{settings: s}, nil
}

func (m *MyLinter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "mylinter",
			Doc:  "just test my linter",
			Run:  m.run,
		},
	}, nil
}

func (m *MyLinter) GetLoadMode() string {
	return register.LoadModeSyntax
}

func (m *MyLinter) run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {

			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			fnName := fn.Name.Name

			returnCount := 0

			if fn.Type.Results != nil {
				for _, field := range fn.Type.Results.List {
					if len(field.Names) == 0 {
						returnCount++
					} else {
						returnCount += len(field.Names)
					}
				}
			}

			if returnCount == 1 {
				retField := fn.Type.Results.List[0]
				ident, ok := retField.Type.(*ast.Ident)
				if ok && ident.Name == "bool" && (len(fnName) > 0 && unicode.IsLower(rune(fnName[0]))) {
					if !strings.HasPrefix(fnName, "is") {
						pass.Reportf(fn.Pos(), "private function '%s' have bool value must start with 'is'", fnName)
					}
				}
			}

			if len(fnName) > 0 && unicode.IsLower(rune(fnName[0])) {
				if !strings.HasPrefix(fnName, "handle") && returnCount > 1 {
					pass.Reportf(fn.Pos(), "private function '%s' have more than one value return must start with 'handle'", fnName)
				}
			}

			return true
		})
	}
	return nil, nil
}
