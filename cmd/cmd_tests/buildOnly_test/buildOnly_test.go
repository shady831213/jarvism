package buildOnly_test

import (
	"github.com/shady831213/jarvism/cmd"
	"github.com/shady831213/jarvism/core"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestBuildOnly(t *testing.T) {
	os.Args = []string{"", "run_build", "build1", "-max_job", "10"}
	if err := cmd.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.PkgPath(), "cmd", "cmd_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
