package core

type Runner interface {
	Name() string
	PrepareBuild(*AstBuild, func(string, ...string) error) *JVSTestResult
	Build(*AstBuild, func(string, ...string) error) *JVSTestResult
	PrepareTest(*AstTestCase, func(string, ...string) error) *JVSTestResult
	RunTest(*AstTestCase, func(string, ...string) error) *JVSTestResult
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
