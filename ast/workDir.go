package ast

import (
	"os"
	"path/filepath"
)

var workDir string

func SetWorkDir(path string) error {
	_workDir, err := filepath.Abs(os.ExpandEnv(path))
	if err != nil {
		return err
	}
	//check path
	if _, err := os.Stat(_workDir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		//create work dir
		if err := os.Mkdir(_workDir, os.ModePerm); err != nil {
			return err
		}
	}

	workDir = _workDir
	return nil
}

func GetWorkDir() string {
	return workDir
}
