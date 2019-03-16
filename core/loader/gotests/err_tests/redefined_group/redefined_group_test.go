package redefined_build

import (
	"github.com/shady831213/jarvism/core/loader"
	"os"
	"strings"
	"testing"
)

func TestRedefineGroup(t *testing.T) {
	err := loader.Load("testFiles/jarvism_cfg")
	if err == nil || !strings.Contains(err.Error(), "group conflict") {
		t.Error("expect group conflict err but get", err.Error())
		t.FailNow()
	}
}

func init() {
	os.Setenv("JVS_PRJ_HOME", "testFiles")
}
