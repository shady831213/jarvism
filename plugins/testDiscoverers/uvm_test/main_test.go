package main_test

import (
	"fmt"
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/loader"
	"os"
	"path"
	"testing"
)

func TestUvmDiscoverer(t *testing.T) {
	cfg, err := loader.Lex("testFiles/test_discover.yaml")
	if err != nil {
		t.Error(err)
	}
	err = loader.Parse(cfg)
	if err != nil {
		t.Error(err)
	}
	build1 := loader.GetJvsAstRoot().GetBuild("build1")
	build2 := loader.GetJvsAstRoot().GetBuild("build2")
	compare(t, "discoverer of build1 name", "uvm_test", build1.GetTestDiscoverer().Name())
	compare(t, "testDir of build1 name", path.Join(core.GetPrjHome(), "build1_testcases"), build1.GetTestDiscoverer().TestDir())
	compare(t, "testList of build1 name", fmt.Sprint([]string{"test2"}), fmt.Sprint(build1.GetTestDiscoverer().TestList()))
	compare(t, "discoverer of build2 name", "uvm_test", build2.GetTestDiscoverer().Name())
	compare(t, "testDir of build2 name", path.Join(core.GetPrjHome(), "testcases"), build2.GetTestDiscoverer().TestDir())
	compare(t, "testList of build2 name", fmt.Sprint([]string{"test1"}), fmt.Sprint(build2.GetTestDiscoverer().TestList()))
}

func compare(t *testing.T, fields, exp, res string) {
	if exp != res {
		t.Error(fields + " expect " + exp + " but get " + res + "!")
	}
}

func init() {
	os.Setenv("JVS_PRJ_HOME", path.Join(core.TestDiscoverersPath(), "uvm_test", "testFiles"))
}
