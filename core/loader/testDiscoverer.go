package loader

import "github.com/shady831213/jarvism/core/plugin"

type TestDiscoverer interface {
	LoderPlugin
	TestDir() string
	TestCmd() string
	TestList() []string
	IsValidTest(string) bool
	TestFileList() []string
}

func RegisterTestDiscoverer(c func() plugin.Plugin) {
	plugin.RegisterPlugin(plugin.JVSTestDiscovererPlugin, c)
}
