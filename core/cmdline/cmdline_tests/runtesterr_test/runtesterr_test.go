package runtesterr_test

import (
	"fmt"
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/cmdline"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRunTestErr(t *testing.T) {
	os.Args = []string{"", "-test", "test1", "-repeat", "10"}
	err := cmdline.Run()
	if err == nil {
		t.Error("should be fail!")
		t.FailNow()
	}
	fmt.Println(err)
}

func init() {
	abs, _ := filepath.Abs(path.Join(core.CorePath(), "cmdline", "cmdline_tests", "testFiles"))
	os.Setenv("JVS_PRJ_HOME", abs)
}
