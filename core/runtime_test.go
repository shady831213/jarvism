package core_test

import (
	"github.com/shady831213/jarvisSim/core"
	"io"
	"math/rand"
	"os/exec"
	"testing"
	"time"
)

type testRunner struct {
}

func (r *testRunner) Name() string {
	return "testRunner"
}

func (r *testRunner) PrepareBuild(build *core.AstBuild, cmdStdout *io.Writer) error {
	cmd := exec.Command("echo", " prepare build ", build.Name)
	cmd.Stdout = *cmdStdout
	time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	return cmd.Run()
}

func (r *testRunner) Build(build *core.AstBuild, cmdStdout *io.Writer) error {
	cmd := exec.Command("echo", " build build ", build.Name)
	cmd.Stdout = *cmdStdout
	time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	return cmd.Run()
}

func (r *testRunner) PrepareTest(testCase *core.AstTestCase, cmdStdout *io.Writer) *core.JVSTestResult {
	cmd := exec.Command("echo", " prepare test ", testCase.Name)
	cmd.Stdout = *cmdStdout
	time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	if err := cmd.Run(); err != nil {
		return core.JVSTestResultFail(err.Error())
	}
	return core.JVSTestResultPass("")
}

func (r *testRunner) RunTest(testCase *core.AstTestCase, cmdStdout *io.Writer) *core.JVSTestResult {
	cmd := exec.Command("echo", " run test ", testCase.Name)
	cmd.Stdout = *cmdStdout
	time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	if err := cmd.Run(); err != nil {
		return core.JVSTestResultFail(err.Error())
	}
	return core.JVSTestResultPass("")
}

func TestGroup(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunGroup(core.GetJvsAstRoot().GetGroup("group1"), nil); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := core.RunGroup(core.GetJvsAstRoot().GetGroup("group2"), []string{}); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := core.RunGroup(core.GetJvsAstRoot().GetGroup("group3"), []string{}); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestSingleTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunTest("test1", "build1", []string{"-seed 1"}); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestSingleRepeatTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunTest("test1", "build1", []string{"-repeat 10"}); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestRunOnlyBuild(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunOnlyBuild("build1", []string{"-test_phase jarvis"}); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
