package core


type Runner interface {
	Name() string
	Prepare() error
	Compile() error
	RunTest(*AstTestCase) error
	RunGroup(*astGroup) error
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
