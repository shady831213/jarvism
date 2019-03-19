package loader

import "github.com/shady831213/jarvism/core/plugin"

//test discoverer interface
//
//detect valid tests with corresponding build
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
