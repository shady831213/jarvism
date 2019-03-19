package runtest_test

import (
	"github.com/shady831213/jarvism/cmd"
	"github.com/shady831213/jarvism/core"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRunTest(t *testing.T) {
	os.Args = []string{"", "run_test", "build1", "test1", "-repeat", "10"}
	if err := cmd.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.PkgPath(), "cmd", "cmd_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
