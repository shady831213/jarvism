package ast

type TestDiscoverer interface {
	pluginOpts
	TestDir() string
	TestCmd() string
	TestList() []string
	IsValidTest(string) bool
	TestFileList() []string
}

var validTestDiscoverers = make(map[string]func() TestDiscoverer)

func GetTestDiscoverer(key string) TestDiscoverer {
	if v, ok := validTestDiscoverers[key]; ok {
		return v()
	}
	return nil
}

func RegisterTestDiscoverer(d func() TestDiscoverer) {
	inst := d()
	if _, ok := validTestDiscoverers[inst.Name()]; ok {
		panic("testDiscoverer " + inst.Name() + " has been registered!")
	}
	validTestDiscoverers[inst.Name()] = d
}
