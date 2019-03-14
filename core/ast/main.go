package ast

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/errors"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

func Lex(configFile string) (map[interface{}]interface{}, error) {
	reader, err := os.Open(configFile)
	defer reader.Close()
	if err != nil {
		return nil, errors.JVSAstLexError("", err.Error())
	}
	cfg := make(map[interface{}]interface{})
	if err := yaml.NewDecoder(reader).Decode(&cfg); err != nil {
		return nil, errors.JVSAstLexError("", err.Error())
	}
	return cfg, nil
}

func Parse(cfg map[interface{}]interface{}) error {
	if err := core.CheckEnv(); err != nil {
		return err
	}
	if err := AstParse(jvsAstRoot, cfg); err != nil {
		return err
	}
	if err := jvsAstRoot.Link(); err != nil {
		return err
	}
	return nil
}

func LoadBuildInOptions(configFile string) error {
	cfg, err := Lex(configFile)
	if err != nil {
		panic(err)
	}
	if err := CfgToAstItemRequired(cfg, "options", func(item interface{}) *errors.JVSAstError {
		for name, option := range item.(map[interface{}]interface{}) {
			jvsAstRoot.Options[name.(string)] = newAstOption(name.(string))
			if err := AstParse(jvsAstRoot.Options[name.(string)], option.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func init() {
	if err := LoadBuildInOptions(path.Join(core.BuildInOptionPath(), "global_options.yaml")); err != nil {
		panic("Error in loading " + path.Join(core.BuildInOptionPath(), "global_options.yaml") + ":" + err.Error())
	}
}
