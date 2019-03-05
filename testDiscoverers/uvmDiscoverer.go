package testDiscoverers

import (
	"github.com/shady831213/jarvisSim"
	"github.com/shady831213/jarvisSim/parser"
	"github.com/shady831213/jarvisSim/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type uvmDiscoverer struct {
	testDir string
	tests   map[string]interface{}
}

func (d *uvmDiscoverer) Parse(cfg map[interface{}]interface{}) error {
	//AstParse tests
	if err := parser.CfgToastItemOptional(cfg, "test_dir", func(item interface{}) error {
		testDir, err := filepath.Abs(os.ExpandEnv(item.(string)))
		if err != nil {
			return err
		}
		//check path
		if _, err := os.Stat(testDir); err != nil {
			return err
		}
		d.testDir = testDir
		return nil
	}); err != nil {
		return parser.AstError("test_dir of "+d.Name(), err)
	}
	//use default
	if d.testDir == "" {
		d.testDir, _ = filepath.Abs(path.Join(jarivsSim.GetPrjHome(), "testcases"))
	}
	return nil
}

func (d *uvmDiscoverer) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"test_dir": nil}
	if !parser.CheckKeyWord(s, keywords) {
		return false, utils.KeyOfStringMap(keywords), "Error in " + d.Name() + ":"
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
	_, ok := d.tests[test]
	return ok
}

func (d *uvmDiscoverer) Clone() parser.TestDiscoverer {
	inst := new(uvmDiscoverer)
	inst.testDir = d.testDir
	if d.tests != nil {
		inst.tests = make(map[string]interface{})
		for k, v := range d.tests {
			inst.tests[k] = v
		}
	}
	return inst
}

func newUvmDiscoverer() parser.TestDiscoverer {
	inst := new(uvmDiscoverer)
	return inst
}

func init() {
	parser.RegisterTestDiscoverer(newUvmDiscoverer())
}
