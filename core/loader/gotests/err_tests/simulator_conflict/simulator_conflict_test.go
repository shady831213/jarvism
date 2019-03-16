package redefined_build

import (
	"github.com/shady831213/jarvism/core/loader"
	"os"
	"strings"
	"testing"
)

func TestSimulatorConflict(t *testing.T) {
	err := loader.Load("testFiles/jarvism_cfg")
	if err == nil || !strings.Contains(err.Error(), "simulator conflict") {
		t.Error("expect simulator conflict err but get", err.Error())
		t.FailNow()
	}
}

func init() {
	os.Setenv("JVS_PRJ_HOME", "testFiles")
}
