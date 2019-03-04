package parser

type Simulator interface {
	Name() string
	BuildInOptionFile() string
	SimCmd() string
	CompileCmd() string
	SeedOption() string
}

var simulator Simulator
var validSimulators = make(map[string]Simulator)

func setSimulator(s Simulator) {
	simulator = s
}

func RegisterSimulator(s Simulator) {
	if _, ok := validSimulators[s.Name()]; ok {
		panic("simulator " + s.Name() + " has been registered!")
	}
	validSimulators[s.Name()] = s
}

func GetSimulator() Simulator {
	return simulator
}
