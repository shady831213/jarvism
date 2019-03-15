package main

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type uvmDiscoverer struct {
	testDir string
	tests   map[string]interface{}
}

func (d *uvmDiscoverer) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	//AstParse test_dir
	if err := loader.CfgToAstItemOptional(cfg, "test_dir", func(item interface{}) *errors.JVSAstError {
		testDir, err := filepath.Abs(os.ExpandEnv(item.(string)))
		if err != nil {
			return errors.JVSAstParseError(d.Name(), err.Error())
		}
		//check path
		if _, err := os.Stat(testDir); err != nil {
			return errors.JVSAstParseError(d.Name(), err.Error())
		}
		d.testDir = testDir
		return nil
	}); err != nil {
		return errors.JVSAstParseError("test_dir of "+d.Name(), err.Error())
	}
	//use default
	if d.testDir == "" {
		d.testDir, _ = filepath.Abs(path.Join(core.GetPrjHome(), "testcases"))
	}
	return nil
}

func (d *uvmDiscoverer) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey("test_dir")
	if !loader.CheckKeyWord(s, keywords) {
		return false, keywords, "Error in " + d.Name() + ":"
	}
	return true, nil, ""
}

func (d *uvmDiscoverer) Name() string {
	return "uvm_test"
}

func (d *uvmDiscoverer) TestDir() string {
	return d.testDir
}

func (d *uvmDiscoverer) TestCmd() string {
	return "+UVM_TESTNAME="
}

func (d *uvmDiscoverer) TestList() []string {
	d.discoverTest()
	return utils.KeyOfStringMap(d.tests)
}

func (d *uvmDiscoverer) TestFileList() []string {
	d.discoverTest()
	values := utils.ValueOfStringMap(d.tests)
	paths := make([]string, len(values))
	for i, v := range values {
		paths[i] = v.(string)
	}
	return paths
}

func (d *uvmDiscoverer) discoverTest() {
	if d.tests == nil {
		d.tests = make(map[string]interface{})
		err := filepath.Walk(d.testDir, d.filter)
		if err != nil {
			panic("Error in test discover :" + err.Error())
		}
	}
}

func (d *uvmDiscoverer) filter(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	if filepath.Ext(path) != ".sv" {
		return nil
	}
	fileName := strings.TrimSuffix(filepath.Base(path), ".sv")
	if fileName == filepath.Base(filepath.Dir(path)) {
		d.tests[fileName] = path
	}
	return nil
}

func (d *uvmDiscoverer) IsValidTest(test string) bool {
	d.discoverTest()
	_, ok := d.tests[test]
	return ok
}

func newUvmDiscoverer() loader.TestDiscoverer {
	inst := new(uvmDiscoverer)
	return inst
}

func init() {
	loader.RegisterTestDiscoverer(newUvmDiscoverer)
}
