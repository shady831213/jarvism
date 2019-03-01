package parser

type vcs struct {
	Simulator
}

func newVcs() *Simulator {
	inst := new(vcs)
	inst.Name = "vcs"
	inst.BuildInOptionFile = "buildInOptions/vcs_options.yaml"
	inst.CompileCmd = "vcs"
	inst.SimCmd = "simv"
	inst.SeedOption = "+ntb_random_seed="
	return &inst.Simulator
}

func init() {
	RegisterSimulator(newVcs())
}
