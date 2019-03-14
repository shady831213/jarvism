package main

import (
	"github.com/shady831213/jarvism/core/ast"
	jvsErrors "github.com/shady831213/jarvism/core/errors"
	"regexp"
)

type testChecker struct {
	ast.CheckerBase
}

func (c *testChecker) Name() string {
	return "testChecker"
}


func newTestChecker() ast.Checker {
	inst := new(testChecker)
	inst.Init()

	//UVM ERROR and FATAL
	inst.AddPats(jvsErrors.JVSRuntimeFail, false, regexp.MustCompile(`^.*UVM_((ERROR)|(FATAL)) .*\@.*:`))
	//jvsErrors
	inst.AddPats(jvsErrors.JVSRuntimeFail, false, regexp.MustCompile(`^Error((.+:)|(-\[.*\]))`))

	//UVM Warning
	inst.AddPats(jvsErrors.JVSRuntimeWarning, false, regexp.MustCompile(`^.*UVM_WARNING .*\@.*:`))
	//Timing violation
	inst.AddPats(jvsErrors.JVSRuntimeWarning, false, regexp.MustCompile(`.*Timing violation.*`))
	return inst
}

func init() {
	ast.RegisterChecker(newTestChecker)
}
