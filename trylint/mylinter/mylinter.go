package mylinter

import (
	"go/ast"

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
			pass.Reportf(f.Pos(), "LINTER IS LIVE: checking file %s", f.Name.Name)
			return true
		})
	}
	return nil, nil
}
