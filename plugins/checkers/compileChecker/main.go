/*
compile checker implementation

"Error:" will be detected
*/

package main

import (
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/plugin"
	"regexp"
)

type compileChecker struct {
	loader.CheckerBase
}

func newCompileChecker() plugin.Plugin {
	inst := new(compileChecker)
	inst.Init("compileChecker")
	//Errors
	inst.AddPats(errors.JVSRuntimeFail, false, regexp.MustCompile(`^Error((.+:)|(-\[.*\]))`))
	return inst
}

func init() {
	loader.RegisterChecker(newCompileChecker)
}
