package main

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

func (r *hostRunner) PrepareBuild(build *ast.AstBuild, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	_, buildName := parseBuildName(build.Name)
	buildDir := path.Join(r.BuildsRoot(), buildName)
	//create build dir
	if err := os.MkdirAll(buildDir, os.ModePerm); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	//create test.f if need
	filelistItem := build.GetTestDiscoverer().TestFileList()
	filelistCompileOption := ""
	hasFileList := filelistItem != nil && len(filelistItem) > 0
	if hasFileList {
		filelistItem = append(filelistItem, build.GetTestDiscoverer().TestDir())
		fileListContent, err := ast.GetSimulator().GetFileList(filelistItem...)
		if err != nil {
			return errors.JVSRuntimeResultFail(err.Error())
		}
		if err := utils.WriteNewFile(path.Join(buildDir, "test.f"), fileListContent); err != nil {
			return errors.JVSRuntimeResultFail(err.Error())
		}
		filelistCompileOption = " -f test.f"
	}

	//create pre_compile,compile, and post_compile script
	if err := utils.WriteNewFile(path.Join(buildDir, "pre_compile.sh"), build.PreCompileAction()); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(buildDir, "compile.sh"), ast.GetSimulator().CompileCmd()+" "+build.CompileOption()+filelistCompileOption); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(buildDir, "post_compile.sh"), build.PostCompileAction()); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(buildDir, "run_compile.sh"), strings.Join([]string{"./pre_compile.sh", bashExitGlue(), "./compile.sh", bashExitGlue(), "./post_compile.sh"}, "\n")); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func (r *hostRunner) Build(build *ast.AstBuild, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	_, buildName := parseBuildName(build.Name)
	buildDir := path.Join(r.BuildsRoot(), buildName)
	//create log file
	logFile, err := os.Create(path.Join(buildDir, buildName+".log"))
	defer logFile.Close()
	if err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	attr := ast.CmdAttr{LogFile: logFile,
		SetAttr: func(cmd *exec.Cmd) error {
			cmd.Dir = path.Join(r.BuildsRoot(), buildName)
			return nil
		}}
	if err := cmdRunner(&attr, "bash", "run_compile.sh"); err != nil {
		return errors.JVSRuntimeResultFail(err.Error() + "\n" + "path:" + buildDir)
	}
	return errors.JVSRuntimeResultPass("path:" + buildDir)
}

func (r *hostRunner) PrepareTest(testCase *ast.AstTestCase, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	_, buildName, testName, seed, groupsName := parseTestName(testCase.Name)
	testDir := path.Join(r.TestsRoot(), path.Join(groupsName...), testName, seed)
	buildDir := path.Join(r.BuildsRoot(), buildName)
	//create test dir
	if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	//link build dir
	if err := os.Symlink(buildDir, path.Join(testDir, buildName)); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	//create pre_sim, sim and post_sim script
	if err := utils.WriteNewFile(path.Join(testDir, "pre_sim.sh"), testCase.PreSimAction()); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(testDir, "sim.sh"), path.Join(buildName, ast.GetSimulator().SimCmd())+" "+testCase.SimOption()+" "+testCase.Build.GetTestDiscoverer().TestCmd()+testName); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(testDir, "post_sim.sh"), testCase.PostSimAction()); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	if err := utils.WriteNewFile(path.Join(testDir, "run_sim.sh"), strings.Join([]string{"./pre_sim.sh", bashExitGlue(), "./sim.sh", bashExitGlue(), "./post_sim.sh"}, "\n")); err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	return errors.JVSRuntimeResultPass("")
}

func (r *hostRunner) RunTest(testCase *ast.AstTestCase, cmdRunner func(*ast.CmdAttr, string, ...string) error) *errors.JVSRuntimeResult {
	_, _, testName, seed, groupsName := parseTestName(testCase.Name)
	testDir := path.Join(r.TestsRoot(), path.Join(groupsName...), testName, seed)
	//create log file
	logFile, err := os.Create(path.Join(testDir, testName+"__"+seed+".log"))
	defer logFile.Close()
	if err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	attr := ast.CmdAttr{LogFile: logFile,
		SetAttr: func(cmd *exec.Cmd) error {
			cmd.Dir = testDir
			return nil
		}}
	if err := cmdRunner(&attr, "bash", "run_sim.sh"); err != nil {
		return errors.JVSRuntimeResultFail(err.Error() + "\n" + "path:" + testDir)
	}
	return errors.JVSRuntimeResultPass("path:" + testDir)
}

func init() {
	ast.RegisterRunner(newHostRunner())
}
