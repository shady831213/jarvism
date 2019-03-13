package main

import (
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
	"math/rand"
	"time"
)

type testRunner struct {
}

func (r *testRunner) Name() string {
	return "test"
}

func (r *testRunner) PrepareBuild(build *ast.AstBuild, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner(nil, "echo", " "); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func (r *testRunner) Build(build *ast.AstBuild, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner(nil, "echo", " build build ", build.Name); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func (r *testRunner) PrepareTest(testCase *ast.AstTestCase, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner(nil, "echo", ""); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func (r *testRunner) RunTest(testCase *ast.AstTestCase, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner(nil, "echo", " run test ", testCase.Name); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func init() {
	ast.RegisterRunner(new(testRunner))
}
