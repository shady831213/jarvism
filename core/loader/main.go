package loader

import (
	"github.com/shady831213/jarvism/core"
	"os"
	"path/filepath"
)

func parseFile(configFile string) error {
	cfg, err := Lex(configFile)
	if err != nil {
		return err
	}
	if err := Parse(cfg); err != nil {
		return err
	}
	return nil
}

func Load(config string) error {
	if err := core.CheckEnv(); err != nil {
		return err
	}
	stat, err := os.Stat(os.ExpandEnv(config))
	if err != nil {
		return err
	}
	//is file
	if !stat.IsDir() {
		if err := parseFile(config); err != nil {
			return err
		}
	} else {
		//dir
		if err := filepath.Walk(config, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if info.IsDir() {
				return err
			}
			if filepath.Ext(path) == ".yaml" {
				return parseFile(path)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	if err := Link(); err != nil {
		return err
	}
	return nil
}
