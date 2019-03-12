package runners

import (
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
)

type hostRunner struct {
}

func (r *hostRunner) Name() string {
	return "host"
}

func (r *hostRunner) PrepareBuild(build *ast.AstBuild, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	return errors.JVSTestResultPass("")
}

func (r *hostRunner) Build(build *ast.AstBuild, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	return errors.JVSTestResultPass("")
}

func (r *hostRunner) PrepareTest(testCase *ast.AstTestCase, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	return errors.JVSTestResultPass("")
}

func (r *hostRunner) RunTest(testCase *ast.AstTestCase, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	return errors.JVSTestResultPass("")
}

func init() {
	ast.RegisterRunner(new(hostRunner))
}
