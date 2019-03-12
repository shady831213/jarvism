package runners

import (
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/utils"
	"os"
	"os/exec"
	"path"
	"strings"
)

type hostRunner struct {
}

func newHostRunner() *hostRunner {
	inst := new(hostRunner)
	return inst
}

func (r *hostRunner) Name() string {
	return "host"
}

func (r *hostRunner) BuildsRoot() string {
	return path.Join(ast.GetWorkDir(), "builds")
}

func (r *hostRunner) TestsRoot() string {
	return path.Join(ast.GetWorkDir(), "tests")
}

func (r *hostRunner) GroupsRoot() string {
	return path.Join(ast.GetWorkDir(), "groups")
}

func (r *hostRunner) PrepareBuild(build *ast.AstBuild, cmdRunner func(func(cmd *exec.Cmd) error, string, ...string) error) *errors.JVSRuntimeResult {
	_, buildName := parseBuildName(build.Name)
	buildDir := path.Join(r.BuildsRoot(), buildName)
	//create build dir
	if err := os.MkdirAll(buildDir, os.ModePerm); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	//create pre_compile,compile and post_compile script
	if err := utils.WriteNewFile(path.Join(buildDir, "pre_compile.sh"), build.PreCompileAction()); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(buildDir, "compile.sh"), ast.GetSimulator().CompileCmd()+" "+build.CompileOption()); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(buildDir, "post_compile.sh"), build.PostCompileAction()); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	glue := "EXCODE=$?\nif [ $EXCODE != 0 ]\nthen\nexit $EXCODE\nfi"
	if err := utils.WriteNewFile(path.Join(buildDir, "run_compile.sh"), strings.Join([]string{"./pre_compile.sh", glue, "./compile.sh",glue, "./post_compile.sh"}, "\n")); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func (r *hostRunner) Build(build *ast.AstBuild, cmdRunner func(func(cmd *exec.Cmd) error, string, ...string) error) *errors.JVSRuntimeResult {
	_, buildName := parseBuildName(build.Name)
	if err := cmdRunner(func(cmd *exec.Cmd) error {
		cmd.Dir = path.Join(r.BuildsRoot(), buildName)
		return nil
	}, "bash", "run_compile.sh"); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func (r *hostRunner) PrepareTest(testCase *ast.AstTestCase, cmdRunner func(func(cmd *exec.Cmd) error, string, ...string) error) *errors.JVSRuntimeResult {

	return errors.JVSRuntimeResultPass("")
}

func (r *hostRunner) RunTest(testCase *ast.AstTestCase, cmdRunner func(func(cmd *exec.Cmd) error, string, ...string) error) *errors.JVSRuntimeResult {
	return errors.JVSRuntimeResultPass("")
}

func init() {
	ast.RegisterRunner(newHostRunner())
}
