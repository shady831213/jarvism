package redefined_build

import (
	"github.com/shady831213/jarvism/core/loader"
	"os"
	"strings"
	"testing"
)

func TestRedefineOption(t *testing.T) {
	err := loader.Load("testFiles/jarvism_cfg")
	if err == nil || !strings.Contains(err.Error(), "option conflict") {
		t.Error("expect option conflict err but get", err.Error())
		t.FailNow()
	}
}

func init() {
	os.Setenv("JVS_PRJ_HOME", "testFiles")
}
