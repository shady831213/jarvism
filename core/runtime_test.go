package core_test

import (
	"github.com/shady831213/jarvisSim/core"
	"io"
	"os/exec"
	"testing"
)

type testRunner struct {
}

func (r *testRunner) Name() string {
	return "testRunner"
}

func (r *testRunner) PrepareBuild(build *core.AstBuild, cmdStdout *io.Writer) error {
	cmd := exec.Command("echo", " prepare build ", build.Name)
	cmd.Stdout = *cmdStdout
	return cmd.Run()
}

func (r *testRunner) Build(build *core.AstBuild, cmdStdout *io.Writer) error {
	cmd := exec.Command("echo", " build build ", build.Name)
	cmd.Stdout = *cmdStdout
	return cmd.Run()
}

func (r *testRunner) PrepareTest(testCase *core.AstTestCase, cmdStdout *io.Writer) error {
	cmd := exec.Command("echo", " prepare test ", testCase.Name)
	cmd.Stdout = *cmdStdout
	return cmd.Run()
}

func (r *testRunner) RunTest(testCase *core.AstTestCase, cmdStdout *io.Writer) error {
	cmd := exec.Command("echo", " run test ", testCase.Name)
	cmd.Stdout = *cmdStdout
	return cmd.Run()
}

func TestGroupTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	core.RunGroup(core.GetJvsAstRoot().GetGroup("group1"))
	core.RunGroup(core.GetJvsAstRoot().GetGroup("group2"))
}

func TestSingleTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunTest("test1", "build1", []string{"-seed 1"}); err != nil {
		t.Error(err)
	}
}

func TestSingleRepeatTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunTest("test1", "build1", []string{"-repeat 10"}); err != nil {
		t.Error(err)
	}
}
