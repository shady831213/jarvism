package runtest_test

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/cmdline"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRunTest(t *testing.T) {
	os.Args = []string{"", "-build", "build1", "-test", "test1", "-repeat", "10"}
	if err := cmdline.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.CorePath(), "cmdline", "cmdline_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
