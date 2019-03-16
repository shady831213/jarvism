package main_test

import (
	"flag"
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/runtime"
	"os"
	"path"
	"testing"
)

var keepResult bool

func tearDown() {
	if !keepResult {
		os.RemoveAll(core.GetWorkDir())
		os.MkdirAll(core.GetWorkDir(), os.ModePerm)
	}
}

func TestHostRunnerBuildFail(t *testing.T) {
	if vcs := os.Getenv("VCS_HOME"); vcs != "" {
		defer tearDown()
		if err := runtime.RunOnlyBuild("build2", nil, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimeFail] != 1 {
			t.Error("expect build fail but it is not!")
			t.FailNow()
		}
	}
}

func TestHostRunnerBuildOnlyAndSimOnly(t *testing.T) {
	if vcs := os.Getenv("VCS_HOME"); vcs != "" {
		defer tearDown()
		if err := runtime.RunOnlyBuild("build1", nil, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 1 {
			t.Error("expect build pass but it is not!")
			t.FailNow()
		}
		//single test
		if err := runtime.RunTest("test1", "build1", []string{"-sim_only"}, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 0 {
			t.Error("expect no build!")
			t.FailNow()
		}
		if runtime.GetTestStatus().Cnts[errors.JVSRuntimePass] != 1 &&
			runtime.GetTestStatus().Cnts[errors.JVSRuntimeFail] != 1 &&
			runtime.GetTestStatus().Cnts[errors.JVSRuntimeWarning] != 1 {
			t.Error("expect test done 1 but it is not!")
			t.FailNow()
		}
	}
}

func TestHostRunnerSim(t *testing.T) {
	if vcs := os.Getenv("VCS_HOME"); vcs != "" {
		defer tearDown()
		//repeat test
		if err := runtime.RunTest("test1", "build1", []string{"-repeat 10"}, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 1 {
			t.Error("expect build pass but it is not!")
			t.FailNow()
		}
		if runtime.GetTestStatus().Cnts[errors.JVSRuntimePass]+runtime.GetTestStatus().Cnts[errors.JVSRuntimeFail]+runtime.GetTestStatus().Cnts[errors.JVSRuntimeWarning] != 10 {
			t.Error("expect test done 10 but it is not!")
			t.FailNow()
		}
		//single test
		if err := runtime.RunTest("test1", "build1", []string{"-sim_only"}, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 0 {
			t.Error("expect no build!")
			t.FailNow()
		}
		if runtime.GetTestStatus().Cnts[errors.JVSRuntimePass] != 1 &&
			runtime.GetTestStatus().Cnts[errors.JVSRuntimeFail] != 1 &&
			runtime.GetTestStatus().Cnts[errors.JVSRuntimeWarning] != 1 {
			t.Error("expect test done 1 but it is not!")
			t.FailNow()
		}
	}
}

func TestHostRunnerGroupSim(t *testing.T) {
	if vcs := os.Getenv("VCS_HOME"); vcs != "" {
		defer tearDown()
		//repeat test
		if err := runtime.RunGroup("group1", []string{"-repeat 10"}, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 2 {
			t.Error("expect build pass 2 but it is not!")
			t.FailNow()
		}
		if runtime.GetTestStatus().Cnts[errors.JVSRuntimePass]+runtime.GetTestStatus().Cnts[errors.JVSRuntimeFail]+runtime.GetTestStatus().Cnts[errors.JVSRuntimeWarning] != 10 {
			t.Error("expect test done 10 but it is not!")
			t.FailNow()
		}
	}
}

func init() {
	os.Setenv("JVS_PRJ_HOME", path.Join(core.RunnersPath(), "host", "testFiles"))
	err := loader.Load("testFiles/runner.yaml")
	if err != nil {
		panic(err)
	}
	flag.BoolVar(&keepResult, "keep", false, "keep test result")
	flag.Parse()
}
