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
	return "testRunner"
}

func (r *testRunner) PrepareBuild(build *ast.AstBuild, cmdRunner ast.CmdRunner) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	return cmdRunner(nil, "echo", " ")
}

func (r *testRunner) Build(build *ast.AstBuild, cmdRunner ast.CmdRunner) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	echoError := rand.Intn(100)
	if echoError < 20 {
		return cmdRunner(nil, "echo", " Error here ", build.Name)
	} else {
		return cmdRunner(nil, "echo", " Pass here ", build.Name)
	}
}

func (r *testRunner) PrepareTest(testCase *ast.AstTestCase, cmdRunner ast.CmdRunner) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	return cmdRunner(nil, "echo", "")
}

func (r *testRunner) RunTest(testCase *ast.AstTestCase, cmdRunner ast.CmdRunner) *errors.JVSRuntimeResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	return cmdRunner(nil, "echo", "UVM_WARNING @abc : ", testCase.Name)
}

func init() {
	ast.RegisterRunner(new(testRunner))
}
