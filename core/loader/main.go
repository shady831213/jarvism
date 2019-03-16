package loader

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/errors"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

func Lex(configFile string) (cfg map[interface{}]interface{}, err error) {
	if filepath.Ext(configFile) != ".yaml" {
		return nil, errors.JVSAstLexError(configFile, "file ext must be .yaml!")
	}
	reader, err := os.Open(configFile)
	defer func() {
		err = reader.Close()
	}()
	if err != nil {
		return nil, errors.JVSAstLexError("", err.Error())
	}
	cfg = make(map[interface{}]interface{})
	if err := yaml.NewDecoder(reader).Decode(&cfg); err != nil {
		return nil, errors.JVSAstLexError("", err.Error())
	}
	return cfg, nil
}

func Parse(cfg map[interface{}]interface{}) error {
	if err := AstParse(jvsAstRoot, cfg); err != nil {
		return err
	}
	return nil
}

func Link() error {
	if err := jvsAstRoot.Link(); err != nil {
		return err
	}
	return nil
}

func ParseFile(configFile string) error {
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
		if err := ParseFile(config); err != nil {
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
				return ParseFile(path)
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
