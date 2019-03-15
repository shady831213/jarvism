package ast

import (
	"github.com/shady831213/jarvism/core/errors"
	"io"
	"os/exec"
	"strings"
)

type CmdAttr struct {
	WriteClosers []io.WriteCloser
	SetAttr      func(*exec.Cmd) error
}

type CmdRunner func(attr *CmdAttr, name string, arg ...string) *errors.JVSRuntimeResult

type Runner interface {
	pluginOpts
	PrepareBuild(*AstBuild, CmdRunner) *errors.JVSRuntimeResult
	Build(*AstBuild, CmdRunner) *errors.JVSRuntimeResult
	PrepareTest(*AstTestCase, CmdRunner) *errors.JVSRuntimeResult
	RunTest(*AstTestCase, CmdRunner) *errors.JVSRuntimeResult
}

func ParseBuildName(name string) (jobId, buildName string) {
	s := strings.Split(name, "__")
	return s[0], s[1]
}

func ParseTestName(name string) (jobId, buildName, testName, seed string, groupsName []string) {
	s := strings.Split(name, "__")
	jobId = s[0]
	buildName = s[1]
	groupsName = s[2 : len(s)-2]
	testName = s[len(s)-2]
	seed = s[len(s)-1]
	return
}

var validRunners = make(map[string]Runner)

func getRunner(key string) Runner {
	if v, ok := validRunners[key]; ok {
		return v
	}
	return nil
}

func RegisterRunner(r Runner) {
	if _, ok := validRunners[r.Name()]; ok {
		panic("runner " + r.Name() + " has been registered!")
	}
	validRunners[r.Name()] = r
}

func GetCurRunner() Runner {
	return jvsAstRoot.env.runner.plugin.(Runner)
}
