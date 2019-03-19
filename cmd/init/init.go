package init

import (
	"errors"
	"fmt"
	"github.com/shady831213/jarvism/cmd/base"
	"github.com/shady831213/jarvism/core/utils"
	"os"
	"path"
)

var CmdInit = &base.Command{
	UsageLine: "jarvism init [-prj_dir DIR][-work_dir DIR]",
	Short:     "create a jarvism default project",
	Long:
	`-$prj_dir
		-jarvism_cfg
			-jarvism_cfg.yaml
			-jarvism_setup.sh(export $JVS_PRJ_HOME;export $JVS_WORK_DIR)
		-src
		-testcases`,
}

var (
	prjDir  string
	workDir string
)

func init() {
	CmdInit.Run = runInit
	CmdInit.Flag.StringVar(&prjDir, "prj_dir", "", "assign prj dir, default is pwd")
	CmdInit.Flag.StringVar(&workDir, "work_dir", "", "assign work dir, default is $prj_dir/work")
	base.Jarvism.AddCommand(CmdInit)
}

func runInit(cmd *base.Command, args []string) error {
	if prjDir == "" {
		prjDir = os.Getenv("PWD")
	} else {
		prjDir = os.ExpandEnv(prjDir)
		if err := os.Mkdir(prjDir, os.ModePerm); err != nil {
			return errors.New(utils.Red(err.Error()))
		}
	}
	//make dirs
	if err := os.Mkdir(path.Join(prjDir, "jarvism_cfg"), os.ModePerm); err != nil {
		return errors.New(utils.Red(err.Error()))
	}
	if err := os.Mkdir(path.Join(prjDir, "src"), os.ModePerm); err != nil {
		return errors.New(utils.Red(err.Error()))
	}
	if err := os.Mkdir(path.Join(prjDir, "testcases"), os.ModePerm); err != nil {
		return errors.New(utils.Red(err.Error()))
	}
	if workDir == "" {
		workDir = path.Join(prjDir, "work")
	}
	workDir = os.ExpandEnv(workDir)
	//files
	setupContent := fmt.Sprintf("#!/bin/bash\nexport JVS_PRJ_HOME=%s\nexprt JVS_WORK_DIR=%s\n", prjDir, workDir)
	if err := utils.WriteNewFile(path.Join(prjDir, "jarvism_cfg", "jarvism_setup.sh"), setupContent); err != nil {
		return errors.New(utils.Red(err.Error()))
	}
	yamlContent :=
		`builds:
			build1:
`
	if err := utils.WriteNewFile(path.Join(prjDir, "jarvism_cfg", "jarvism_cfg.yaml"), yamlContent); err != nil {
		return errors.New(utils.Red(err.Error()))
	}

	return nil
}
