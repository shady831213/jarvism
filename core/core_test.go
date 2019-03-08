package core_test

import (
	"github.com/shady831213/jarvism/core"
	"os"
)

func init() {
	os.Setenv("JVS_PRJ_HOME", "testFiles")
	cfg, err := core.Lex("testFiles/build.yaml")
	if err != nil {
		panic(err)
	}
	err = core.Parse(cfg)
	if err != nil {
		panic(err)
	}
}
