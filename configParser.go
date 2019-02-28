package jarivsSim

import (
	"gopkg.in/yaml.v2"
	"os"
)

func Lex(configFile string) (map[interface{}]interface{}, error) {
	reader, err := os.Open(configFile)
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	cfg := make(map[interface{}]interface{})
	err = yaml.NewDecoder(reader).Decode(&cfg)
	return cfg, err
}

func Parse(cfg map[interface{}]interface{}) error {
	if err := jvsASTRoot.Parse(cfg); err != nil {
		return err
	}
	if err := jvsASTRoot.Link(); err != nil {
		return err
	}
	return nil
}
