package main

import (
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/plugin"
	"github.com/shady831213/jarvism/core/utils"
)

type testSim struct {
}

func newTestSim() plugin.Plugin {
	return new(testSim)
}

func (s *testSim) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	return nil
}

func (s *testSim) KeywordsChecker(key string) (bool, *utils.StringMapSet, string) {
	return true, nil, ""
}

func (s *testSim) Name() string {
	return "testSim"
}

func (s *testSim) BuildInOptionFile() string {
	return ""
}

func (s *testSim) CompileCmd() string {
	return ""
}

func (s *testSim) SimCmd() string {
	return ""
}

func (s *testSim) SeedOption() string {
	return ""
}

func (s *testSim) GetFileList(paths ...string) (string, error) {
	return "", nil
}

func init() {
	loader.RegisterSimulator(newTestSim)
}
