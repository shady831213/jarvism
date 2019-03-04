package simulators

import (
	"github.com/shady831213/jarvisSim"
	"github.com/shady831213/jarvisSim/parser"
	"path"
)

type vcs struct {
}

func (s *vcs) Name() string {
	return "vcs"
}

func (s *vcs) BuildInOptionFile() string {
	return path.Join(jarivsSim.BuildInOptionPath(), "vcs_options.yaml")
}

func (s *vcs) CompileCmd() string {
	return "vcs"
}

func (s *vcs) SimCmd() string {
	return "simv"
}

func (s *vcs) SeedOption() string {
	return "+ntb_random_seed="
}

func newVcs() parser.Simulator {
	inst := new(vcs)
	return inst
}

func init() {
	parser.RegisterSimulator(newVcs())
}
