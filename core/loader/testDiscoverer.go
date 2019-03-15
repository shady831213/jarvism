package loader

type TestDiscoverer interface {
	Plugin
	TestDir() string
	TestCmd() string
	TestList() []string
	IsValidTest(string) bool
	TestFileList() []string
}

func RegisterTestDiscoverer(c func() Plugin) {
	registerPlugin(JVSTestDiscovererPlugin, c)
}
