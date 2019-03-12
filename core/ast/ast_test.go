package ast_test

import (
	"github.com/shady831213/jarvism/core/ast"
	_ "github.com/shady831213/jarvism/runners"
	_ "github.com/shady831213/jarvism/simulators"
	_ "github.com/shady831213/jarvism/testDiscoverers"
	"os"
)

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
