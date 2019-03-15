package ast

import (
	"github.com/shady831213/jarvism/core/errors"
	"io"
	"os/exec"
)

type CmdAttr struct {
	WriteClosers []io.WriteCloser
	SetAttr      func(*exec.Cmd) error
}

type CmdRunner func(attr *CmdAttr, name string, arg ...string) *errors.JVSRuntimeResult

type Runner interface {
	Name() string
	PrepareBuild(*AstBuild, CmdRunner) *errors.JVSRuntimeResult
	Build(*AstBuild, CmdRunner) *errors.JVSRuntimeResult
	PrepareTest(*AstTestCase, CmdRunner) *errors.JVSRuntimeResult
	RunTest(*AstTestCase, CmdRunner) *errors.JVSRuntimeResult
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
