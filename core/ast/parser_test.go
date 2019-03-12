package ast_test

import (
	"github.com/shady831213/jarvism"
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/utils"
	"path"
	"regexp"
	"strings"
	"syscall"
	"testing"
)

func TestLex(t *testing.T) {
	cfg, err := ast.Lex(path.Join(jarivsm.CorePath(), "testFiles/build.yaml"))
	if err != nil {
		t.Error(err)
	}
	t.Log(cfg)
}

func TestParse(t *testing.T) {
	expect := utils.ReadFile(path.Join(jarivsm.CorePath(), "testFiles/build.ast"))
	result := ast.GetJvsAstRoot().GetHierString(0)
	result = dealAstResult(result)

	if result != expect {
		t.Log(ast.GetJvsAstRoot().GetHierString(0))
		utils.WriteNewFile(path.Join(jarivsm.CorePath(), "testFiles/build.ast.result"), result)
		t.Error("not equal! please diff testFiles/build.ast and estFiles/build.ast.result")
		return
	}
	syscall.Unlink(path.Join(jarivsm.CorePath(), "testFiles/build.ast.result"))
}

func dealAstResult(result string) string {
	_result := strings.Replace(result, " ", "", -1)
	//replace rand numbers
	_result = regexp.MustCompile(`(.*__)\d+`).ReplaceAllString(_result, "${1}seed")
	_result = regexp.MustCompile(`\+ntb_random_seed=\d+`).ReplaceAllString(_result, "+ntb_random_seed=seed")
	_result = regexp.MustCompile(`\[[0-9]+\]`).ReplaceAllString(_result, "[seeds]")
	return _result

}
