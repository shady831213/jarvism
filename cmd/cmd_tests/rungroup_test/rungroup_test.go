package rungroup_test

import (
	"github.com/shady831213/jarvism/cmd"
	"github.com/shady831213/jarvism/core"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRunGroup(t *testing.T) {
	os.Args = []string{"", "run_group", "group3", "-sim_only", "-reporter", "junit", "-sim_args", "\"-abc def\""}
	if err := cmd.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.PkgPath(), "cmd", "cmd_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
