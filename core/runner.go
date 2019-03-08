package core

import "io"

type Runner interface {
	Name() string
	PrepareBuild(*AstBuild, *io.Writer) error
	Build(*AstBuild, *io.Writer) error
	PrepareTest(*AstTestCase, *io.Writer) *JVSTestResult
	RunTest(*AstTestCase, *io.Writer) *JVSTestResult
}

var runner Runner
var validRunners = make(map[string]Runner)

//Fix Me: export for test
func SetRunner(r Runner) {
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
