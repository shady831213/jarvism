package options_test

import (
	"github.com/shady831213/jarvism/core/loader"
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
	err := loader.Load("testFiles/build.yaml")
	if err != nil {
		panic(err)
	}
}
