package showPlugins_test

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/cmdline"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestShowPlugins(t *testing.T) {
	os.Args = []string{"", "-show_plugins", "all"}
	if err := cmdline.Run(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.CorePath(), "cmdline", "cmdline_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
