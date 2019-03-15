package loader

type Simulator interface {
	Plugin
	BuildInOptionFile() string
	SimCmd() string
	CompileCmd() string
	SeedOption() string
	GetFileList(...string) (string, error)
}


func RegisterSimulator(c func() Plugin) {
	registerPlugin(JVSSimulatorPlugin, c)
}

func GetCurSimulator() Simulator {
	return jvsAstRoot.env.simulator.plugin.(Simulator)
}
