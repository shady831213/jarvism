package parser_test

import (
	"github.com/shady831213/jarvisSim/parser"
	_ "github.com/shady831213/jarvisSim/simulators"
	"github.com/shady831213/jarvisSim/utils"
	"math/rand"
	"strings"
	"syscall"
	"testing"
)

func TestLex(t *testing.T) {
	cfg, err := parser.Lex("testFiles/build.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Log(cfg)
}

func TestParse(t *testing.T) {
	parser.SetRand(rand.New(rand.NewSource(1)))
	cfg, err := parser.Lex("testFiles/build.yaml")
	if err != nil {
		t.Error(err)
	}
	err = parser.Parse(cfg)
	if err != nil {
		t.Error(err)
	}
	expect := utils.ReadFile("testFiles/build.ast")
	result := parser.GetJvsAstRoot().GetHierString(0)
	result = strings.Replace(result, " ", "", -1)
	if result != expect {
		t.Log(parser.GetJvsAstRoot().GetHierString(0))
		utils.WriteNewFile("testFiles/build.ast.result", result)
		t.Error("not equal! please diff testFiles/build.ast and estFiles/build.ast.result")
		return
	}
	syscall.Unlink("testFiles/build.ast.result")
}
