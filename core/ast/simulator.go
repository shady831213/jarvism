package ast

type Simulator interface {
	pluginOpts
	BuildInOptionFile() string
	SimCmd() string
	CompileCmd() string
	SeedOption() string
	GetFileList(...string) (string, error)
}

var simulator Simulator
var validSimulators = make(map[string]Simulator)

func setSimulator(s Simulator) {
	simulator = s
}

func GetSimulator(key string) Simulator {
	if v, ok := validSimulators[key]; ok {
		return v
	}
	return nil
}

func RegisterSimulator(s Simulator) {
	if _, ok := validSimulators[s.Name()]; ok {
		panic("simulator " + s.Name() + " has been registered!")
	}
	validSimulators[s.Name()] = s
}

func GetCurSimulator() Simulator {
	return simulator
}
