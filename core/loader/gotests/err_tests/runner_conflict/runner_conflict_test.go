package redefined_build

import (
	"github.com/shady831213/jarvism/core/loader"
	"os"
	"strings"
	"testing"
)

func TestRunnerConflict(t *testing.T) {
	err := loader.Load("testFiles/jarvism_cfg")
	if err == nil || !strings.Contains(err.Error(), "runner conflict") {
		t.Error("expect runner conflict err but get", err.Error())
		t.FailNow()
	}
}

func init() {
	os.Setenv("JVS_PRJ_HOME", "testFiles")
}
