package showGroups_test

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/cmdline"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestShowGroups(t *testing.T) {
	os.Args = []string{"", "-show_groups", "-max_job", "10"}
	if err := cmdline.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.CorePath(), "cmdline", "cmdline_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
