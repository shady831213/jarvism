package jarivsSim

import (
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
	t.Log(jvsASTRoot)
}
