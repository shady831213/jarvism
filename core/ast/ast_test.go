package ast_test

import (
	"github.com/shady831213/jarvism"
	"github.com/shady831213/jarvism/core/ast"
	_ "github.com/shady831213/jarvism/simulators"
	_ "github.com/shady831213/jarvism/testDiscoverers"
	"os"
	"path"
)

func init() {
	os.Setenv("JVS_PRJ_HOME", path.Join(jarivsm.CorePath(), "testFiles"))
	cfg, err := ast.Lex(path.Join(jarivsm.CorePath(), "testFiles/build.yaml"))
	if err != nil {
		panic(err)
	}
	err = ast.Parse(cfg)
	if err != nil {
		panic(err)
	}
}
