package parser

type TestDiscoverer interface {
	astParser
	Name() string
	TestDir() string
	TestCmd() string
	TestList() []string
	IsValidTest(string) bool
	Clone() TestDiscoverer
}

var validTestDiscoverers = make(map[string]TestDiscoverer)

func GetTestDiscoverer(key string) TestDiscoverer {
	if v, ok := validTestDiscoverers[key]; ok {
		return v.Clone()
	}
	return nil
}

func RegisterTestDiscoverer(d TestDiscoverer) {
	if _, ok := validTestDiscoverers[d.Name()]; ok {
		panic("testDiscoverer " + d.Name() + " has been registered!")
	}
	validTestDiscoverers[d.Name()] = d
}
