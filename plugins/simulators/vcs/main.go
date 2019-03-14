package main

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/ast"
	"os"
	"path"
	"path/filepath"
)

type vcs struct {
}

func (s *vcs) Name() string {
	return "vcs"
}

func (s *vcs) BuildInOptionFile() string {
	return path.Join(core.SimulatorsPath(), s.Name(), "buildInOptions", "vcs_options.yaml")
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

func (s *vcs) GetFileList(paths ...string) (string, error) {
	fileList := ""
	for _, p := range paths {
		item, err := filepath.Abs(os.ExpandEnv(p))
		if err != nil {
			return "", err
		}
		//check stat
		stat, err := os.Stat(item)
		if err != nil {
			return "", err
		}
		if stat.IsDir() {
			fileList += "+incdir+" + item + "\n"
		} else {
			fileList += item + "\n"
		}
	}
	return fileList, nil
}

func newVcs() ast.Simulator {
	inst := new(vcs)
	return inst
}

func init() {
	ast.RegisterSimulator(newVcs())
}
