/*
sim checker implementation

"Error:" will be detected

"UVM_ERROR:" will be detected

"UVM_FATAL:" will be detected

"UVM_WARNING:" will be detected

"Timing violation:" will be detected
*/

package main

import (
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/plugin"
	"regexp"
)

type testChecker struct {
	loader.CheckerBase
}

func newTestChecker() plugin.Plugin {
	inst := new(testChecker)
	inst.Init("testChecker")

	//UVM ERROR and FATAL
	inst.AddPats(errors.JVSRuntimeFail, false, regexp.MustCompile(`^.*UVM_((ERROR)|(FATAL)) .*\@.*:`))
	//Errors
	inst.AddPats(errors.JVSRuntimeFail, false, regexp.MustCompile(`^Error((.+:)|(-\[.*\]))`))

	//UVM Warning
	inst.AddPats(errors.JVSRuntimeWarning, false, regexp.MustCompile(`^.*UVM_WARNING .*\@.*:`))
	//Timing violation
	inst.AddPats(errors.JVSRuntimeWarning, false, regexp.MustCompile(`.*Timing violation.*`))
	return inst
}

func init() {
	loader.RegisterChecker(newTestChecker)
}
