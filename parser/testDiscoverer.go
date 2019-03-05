package parser

type TestDiscoverer interface {
	astParser
	Name() string
	TestDir() string
	TestCmd() string
	TestList() []string
	IsValidTest(string) bool
}

var validTestDiscoverers = make(map[string]TestDiscoverer)

func RegisterTestDiscoverer(d TestDiscoverer) {
	if _, ok := validTestDiscoverers[d.Name()]; ok {
		panic("testDiscoverer " + d.Name() + " has been registered!")
	}
	validTestDiscoverers[d.Name()] = d
}
