package loader

import "github.com/shady831213/jarvism/core/plugin"

type Simulator interface {
	LoderPlugin
	BuildInOptionFile() string
	SimCmd() string
	CompileCmd() string
	SeedOption() string
	GetFileList(...string) (string, error)
}

func RegisterSimulator(c func() plugin.Plugin) {
	plugin.RegisterPlugin(plugin.JVSSimulatorPlugin, c)
}

func GetCurSimulator() Simulator {
	return jvsAstRoot.env.simulator.plugin.(Simulator)
}
