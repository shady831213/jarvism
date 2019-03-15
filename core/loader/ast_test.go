package loader_test

import (
	"github.com/shady831213/jarvism/core/loader"
	"os"
)

func init() {
	os.Setenv("JVS_PRJ_HOME", "testFiles")
	cfg, err := loader.Lex("testFiles/build.yaml")
	if err != nil {
		panic(err)
	}
	err = loader.Parse(cfg)
	if err != nil {
		panic(err)
	}
}
