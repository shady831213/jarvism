package jarivsSim

import (
	"strings"
	"syscall"
	"testing"
)

func TestLex(t *testing.T) {
	cfg, err := Lex("testFiles/build.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Log(cfg)
}

func TestParse(t *testing.T) {
	cfg, err := Lex("testFiles/build.yaml")
	if err != nil {
		t.Error(err)
	}
	err = Parse(cfg)
	if err != nil {
		t.Error(err)
	}
	expect := readFile("testFiles/build.ast")
	result := jvsAstRoot.GetHierString(0)
	result = strings.Replace(result, " ", "", -1)
	if result != expect {
		t.Log(jvsAstRoot.GetHierString(0))
		writeNewFile("testFiles/build.ast.result", result)
		t.Error("not equal! please diff testFiles/build.ast and estFiles/build.ast.result")
		return
	}
	syscall.Unlink("testFiles/build.ast.result")
}
