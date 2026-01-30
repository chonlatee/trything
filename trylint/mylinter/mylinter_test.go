package mylinter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPluginMylinter(t *testing.T) {
	newPlugin, err := register.GetPlugin("mylinter")
	require.NoError(t, err)

	plugin, err := newPlugin(nil)
	require.NoError(t, err)

	analyzers, err := plugin.BuildAnalyzers()
	require.NoError(t, err)

	analysistest.Run(t, testDataDir(t), analyzers[0], "testlintdata/data")
}

func testDataDir(t *testing.T) string {
	t.Helper()

	_, testFilename, _, ok := runtime.Caller(1)

	if !ok {
		require.Fail(t, "unable go get current test file")
	}

	return filepath.Join(filepath.Dir(testFilename), "testdata")

}
