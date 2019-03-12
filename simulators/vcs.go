package mulators

import (
	"github.com/shady831213/jarvism"
	"github.com/shady831213/jarvism/core/ast"
	"path"
)

type vcs struct {
}

func (s *vcs) Name() string {
	return "vcs"
}

func (s *vcs) BuildInOptionFile() string {
	return path.Join(jarivsm.BuildInOptionPath(), "vcs_options.yaml")
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

func newVcs() ast.Simulator {
	inst := new(vcs)
	return inst
}

func init() {
	ast.RegisterSimulator(newVcs())
}
