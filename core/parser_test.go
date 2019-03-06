package core_test

import (
	"github.com/shady831213/jarvisSim/core"
	_ "github.com/shady831213/jarvisSim/simulators"
	_ "github.com/shady831213/jarvisSim/testDiscoverers"
	"github.com/shady831213/jarvisSim/utils"
	"regexp"
	"strings"
	"syscall"
	"testing"
)

func TestLex(t *testing.T) {
	cfg, err := core.Lex("testFiles/build.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Log(cfg)
}

func TestParse(t *testing.T) {
	expect := utils.ReadFile("testFiles/build.ast")
	result := core.GetJvsAstRoot().GetHierString(0)
	result = dealAstResult(result)

	if result != expect {
		t.Log(core.GetJvsAstRoot().GetHierString(0))
		utils.WriteNewFile("testFiles/build.ast.result", result)
		t.Error("not equal! please diff testFiles/build.ast and estFiles/build.ast.result")
		return
	}
	syscall.Unlink("testFiles/build.ast.result")
}

func dealAstResult(result string) string {
	_result := strings.Replace(result, " ", "", -1)
	//replace rand numbers
	_result = regexp.MustCompile(`(.*__)\d+`).ReplaceAllString(_result, "${1}seed")
	_result = regexp.MustCompile(`\+ntb_random_seed=\d+`).ReplaceAllString(_result, "+ntb_random_seed=seed")
	_result = regexp.MustCompile(`\[[0-9]+\]`).ReplaceAllString(_result, "[seeds]")
	return _result

}
