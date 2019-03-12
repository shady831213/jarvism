package ast

import "github.com/shady831213/jarvism/core/errors"

type Runner interface {
	Name() string
	PrepareBuild(*AstBuild, func(string, ...string) error) *errors.JVSTestResult
	Build(*AstBuild, func(string, ...string) error) *errors.JVSTestResult
	PrepareTest(*AstTestCase, func(string, ...string) error) *errors.JVSTestResult
	RunTest(*AstTestCase, func(string, ...string) error) *errors.JVSTestResult
}

var runner Runner
var validRunners = make(map[string]Runner)

func setRunner(r Runner) {
	runner = r
}

func RegisterRunner(r Runner) {
	if _, ok := validRunners[r.Name()]; ok {
		panic("runner " + r.Name() + " has been registered!")
	}
	validRunners[r.Name()] = r
}

func GetRunner() Runner {
	return runner
}
