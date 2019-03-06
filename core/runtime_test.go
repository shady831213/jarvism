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

func (r *testRunner) PrepareBuild(build *core.AstBuild, cmdStdout io.Writer) error {
	cmd := exec.Command("echo", " prepare build ", build.Name)
	cmd.Stdout = cmdStdout
	return cmd.Run()
}

func (r *testRunner) Build(build *core.AstBuild, cmdStdout io.Writer) error {
	cmd := exec.Command("echo", " build build ", build.Name)
	cmd.Stdout = cmdStdout
	return cmd.Run()
}

func (r *testRunner) PrepareTest(testCase *core.AstTestCase, cmdStdout io.Writer) error {
	cmd := exec.Command("echo", " prepare test ", testCase.Name)
	cmd.Stdout = cmdStdout
	return cmd.Run()
}

func (r *testRunner) RunTest(testCase *core.AstTestCase, cmdStdout io.Writer) error {
	cmd := exec.Command("echo", " run test ", testCase.Name)
	cmd.Stdout = cmdStdout
	return cmd.Run()
}

func TestGroupTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	core.Run(core.GetJvsAstRoot().GetGroup("group1"))
}

func TestSingleTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	g := core.NewAstGroup("test1")
	g.Parse(map[interface{}]interface{}{"build": "build1",
		"tests": map[interface{}]interface{}{"test1": map[interface{}]interface{}{"args": []interface{}{"-seed 1"}}}})
	g.Link()
	core.Run(g)
}
