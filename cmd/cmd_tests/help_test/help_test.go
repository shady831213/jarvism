package help_test

import (
	"github.com/shady831213/jarvism/cmd"
	"github.com/shady831213/jarvism/core"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestHelp(t *testing.T) {
	os.Args = []string{"jarvism", "help"}
	if err := cmd.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	os.Args = []string{"jarvism", "help", "run_test"}
	if err := cmd.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	os.Args = []string{"jarvism", "help", "run_build"}
	if err := cmd.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.PkgPath(), "cmd", "cmd_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
