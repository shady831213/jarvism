package options_test

import (
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/options"
	_ "github.com/shady831213/jarvism/core/runtime"
	"os"
	"testing"
)

func TestOptionUsage(t *testing.T) {
	options.GetJvsOptions().Usage()
}

func init() {
	os.Setenv("JVS_PRJ_HOME", "testFiles")
	cfg, err := ast.Lex("testFiles/build.yaml")
	if err != nil {
		panic(err)
	}
	err = ast.Parse(cfg)
	if err != nil {
		panic(err)
	}
}
