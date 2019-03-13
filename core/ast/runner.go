package ast

import (
	"github.com/shady831213/jarvism/core/errors"
	"os"
	"os/exec"
)

type CmdAttr struct {
	LogFile *os.File
	SetAttr func(*exec.Cmd) error
}

type Runner interface {
	Name() string
	PrepareBuild(*AstBuild, func(*CmdAttr, string, ...string) error) *errors.JVSRuntimeResult
	Build(*AstBuild, func(*CmdAttr, string, ...string) error) *errors.JVSRuntimeResult
	PrepareTest(*AstTestCase, func(*CmdAttr, string, ...string) error) *errors.JVSRuntimeResult
	RunTest(*AstTestCase, func(*CmdAttr, string, ...string) error) *errors.JVSRuntimeResult
}

var runner Runner
var validRunners = make(map[string]Runner)

func setRunner(r Runner) {
	runner = r
}

func RegisterRunner(r Runner) {
	if _, ok := validRunners[r.Name()]; ok {
		panic("runner " + r.Name() + " has been registered!")
	}
	validRunners[r.Name()] = r
}

func GetRunner() Runner {
	return runner
}
