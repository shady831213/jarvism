package main

import (
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
	"regexp"
)

type compileChecker struct {
	ast.CheckerBase
}

func newCompileChecker() ast.Checker {
	inst := new(compileChecker)
	inst.Init("compileChecker")
	//Errors
	inst.AddPats(errors.JVSRuntimeFail, false, regexp.MustCompile(`^Error((.+:)|(-\[.*\]))`))
	return inst
}

func init() {
	ast.RegisterChecker(newCompileChecker)
}
