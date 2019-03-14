package main

import (
	"github.com/shady831213/jarvism/core/ast"
)

type compileChecker struct {
	ast.CheckerBase
}

func (c *compileChecker) Name() string {
	return "compileChecker"
}

func newCompileChecker() ast.Checker {
	inst := new(compileChecker)
	inst.Init()
	return inst
}

func init() {
	ast.RegisterChecker(newCompileChecker)
}
