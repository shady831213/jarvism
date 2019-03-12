package runners_test

import (
	"flag"
	"github.com/shady831213/jarvism"
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/runtime"
	_ "github.com/shady831213/jarvism/simulators"
	_ "github.com/shady831213/jarvism/testDiscoverers"
	"os"
	"path"
	"testing"
)

var keepResult bool

func TestHostRunnerBuildFail(t *testing.T) {
	if vcs := os.Getenv("VCS_HOME"); vcs != "" {
		cfg, err := ast.Lex("testFiles/runner_compile_fail.yaml")
		if err != nil {
			t.Error(err)
		}
		err = ast.Parse(cfg)
		if err != nil {
			t.Error(err)
		}
		if err := runtime.RunOnlyBuild("build1", nil, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimeFail] != 1 {
			t.Error("expect build fail but it is not!")
			t.FailNow()
			return
		}
		if !keepResult {
			os.RemoveAll(ast.GetWorkDir())
		}
	}
}

func TestHostRunnerBuild(t *testing.T) {
	if vcs := os.Getenv("VCS_HOME"); vcs != "" {
		cfg, err := ast.Lex("testFiles/runner.yaml")
		if err != nil {
			t.Error(err)
		}
		err = ast.Parse(cfg)
		if err != nil {
			t.Error(err)
		}
		if err := runtime.RunOnlyBuild("build1", nil, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 1 {
			t.Error("expect build pass but it is not!")
			t.FailNow()
			return
		}
		if !keepResult {
			os.RemoveAll(ast.GetWorkDir())
		}
	}
}

func TestHostRunnerSim(t *testing.T) {
	if vcs := os.Getenv("VCS_HOME"); vcs != "" {
		cfg, err := ast.Lex("testFiles/runner.yaml")
		if err != nil {
			t.Error(err)
		}
		err = ast.Parse(cfg)
		if err != nil {
			t.Error(err)
		}
		//repeat test
		if err := runtime.RunTest("test1", "build1", []string{"-repeat 10"}, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 1 {
			t.Error("expect build pass but it is not!")
			t.FailNow()
			return
		}
		if runtime.GetTestStatus().Cnts[errors.JVSRuntimePass] != 10 {
			t.Error("expect test pass 10 but it is not!")
			t.FailNow()
			return
		}
		//single test
		if err := runtime.RunTest("test1", "build1", []string{"-sim_only"}, nil); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if runtime.GetBuildStatus().Cnts[errors.JVSRuntimePass] != 0 {
			t.Error("expect no build!")
			t.FailNow()
			return
		}
		if runtime.GetTestStatus().Cnts[errors.JVSRuntimePass] != 1 {
			t.Error("expect test pass 1 but it is not!")
			t.FailNow()
			return
		}
		if !keepResult {
			os.RemoveAll(ast.GetWorkDir())
		}
	}
}

func init() {
	os.Setenv("JVS_PRJ_HOME", path.Join(jarivsm.RunnersPath(), "testFiles"))
	flag.BoolVar(&keepResult, "keep", false, "keep test result")
	flag.Parse()
}
