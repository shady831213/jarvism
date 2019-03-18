package main_test

import (
	"flag"
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/runtime"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var keepResult bool

func tearDonw() {
	if !keepResult {
		os.RemoveAll(core.GetWorkDir())
	}
}

func TestGroup(t *testing.T) {
	if err := runtime.RunGroup("group2", []string{"-reporter junit"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	tearDonw()
}

func TestOption(t *testing.T) {
	if err := runtime.RunGroup("group3", []string{"-reporter junit", "-junitNoXMLHeader"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	tearDonw()
}

func init() {
	//build and load plugin
	abs, _ := filepath.Abs("testFiles")
	os.Setenv("JVS_PRJ_HOME", abs)
	err := loader.Load(path.Join(abs, "jarvism_cfg"))
	if err != nil {
		panic(err)
	}
	flag.BoolVar(&keepResult, "keep", false, "keep test result")
	flag.Parse()
}
