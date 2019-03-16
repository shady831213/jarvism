package loader_basic_test

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/utils"
	"os"
	"path"
	"regexp"
	"strings"
	"syscall"
	"testing"
)

func TestLex(t *testing.T) {
	cfg, err := loader.Lex(path.Join(core.CorePath(), "loader", "testFiles/build.yaml"))
	if err != nil {
		t.Error(err)
	}
	t.Log(cfg)
}

func TestParse(t *testing.T) {
	err := loader.Load(path.Join(core.CorePath(), "loader", "testFiles/jarvism_cfg"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expect, _ := utils.ReadFile(path.Join(core.CorePath(), "loader", "testFiles/build.ast"))
	result := loader.GetJvsAstRoot().GetHierString(0)
	result = dealAstResult(result)

	if result != expect {
		t.Log(loader.GetJvsAstRoot().GetHierString(0))
		utils.WriteNewFile(path.Join(core.CorePath(), "loader", "testFiles/build.ast.result"), result)
		t.Error("not equal! please diff testFiles/build.ast and estFiles/build.ast.result")
		return
	}
	syscall.Unlink(path.Join(core.CorePath(), "loader", "testFiles/build.ast.result"))
}

func dealAstResult(result string) string {
	_result := strings.Replace(result, " ", "", -1)
	//replace rand numbers
	_result = regexp.MustCompile(`(.*__)\d+`).ReplaceAllString(_result, "${1}seed")
	_result = regexp.MustCompile(`\+ntb_random_seed=\d+`).ReplaceAllString(_result, "+ntb_random_seed=seed")
	_result = regexp.MustCompile(`\[[0-9]+\]`).ReplaceAllString(_result, "[seeds]")
	return _result

}

func init() {
	os.Setenv("JVS_PRJ_HOME", path.Join(core.CorePath(), "loader", "testFiles"))
}
