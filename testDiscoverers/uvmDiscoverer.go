package testDiscoverers

import (
	"github.com/shady831213/jarvism"
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type uvmDiscoverer struct {
	testDir string
	tests   map[string]interface{}
}

func (d *uvmDiscoverer) Parse(cfg map[interface{}]interface{}) *errors.AstError {
	//AstParse tests
	if err := ast.CfgToAstItemOptional(cfg, "test_dir", func(item interface{}) *errors.AstError {
		testDir, err := filepath.Abs(os.ExpandEnv(item.(string)))
		if err != nil {
			return errors.NewAstParseError(d.Name(), err.Error())
		}
		//check path
		if _, err := os.Stat(testDir); err != nil {
			return errors.NewAstParseError(d.Name(), err.Error())
		}
		d.testDir = testDir
		return nil
	}); err != nil {
		return errors.NewAstParseError("test_dir of "+d.Name(), err.Error())
	}
	//use default
	if d.testDir == "" {
		d.testDir, _ = filepath.Abs(path.Join(jarivsm.GetPrjHome(), "testcases"))
	}
	return nil
}

func (d *uvmDiscoverer) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey("test_dir")
	if !ast.CheckKeyWord(s, keywords) {
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
	if d.tests == nil {
		d.tests = make(map[string]interface{})
	}

	err := filepath.Walk(d.testDir, d.filter)
	if err != nil {
		panic("Error in test discover :" + err.Error())
	}

	return utils.KeyOfStringMap(d.tests)
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
		d.tests[fileName] = nil
	}
	return nil
}

func (d *uvmDiscoverer) IsValidTest(test string) bool {
	d.TestList()
	_, ok := d.tests[test]
	return ok
}

func newUvmDiscoverer() ast.TestDiscoverer {
	inst := new(uvmDiscoverer)
	return inst
}

func init() {
	ast.RegisterTestDiscoverer(newUvmDiscoverer)
}
