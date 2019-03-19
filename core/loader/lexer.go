package loader

import (
	"github.com/shady831213/jarvism/core/errors"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

//Lex a yaml file, output a map tree
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
	curLexFile = configFile
	cfg = make(map[interface{}]interface{})
	if err := yaml.NewDecoder(reader).Decode(&cfg); err != nil {
		return nil, errors.JVSAstLexError("", err.Error())
	}
	return cfg, nil
}

var curLexFile string
