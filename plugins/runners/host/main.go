package main

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/utils"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func bashExitGlue() string {
	return "EXCODE=$?\nif [ $EXCODE != 0 ]\nthen\nexit $EXCODE\nfi"
}

type hostRunner struct {
}

func newHostRunner() loader.Plugin {
	inst := new(hostRunner)
	return inst
}

func (r *hostRunner) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	return nil
}

func (r *hostRunner) KeywordsChecker(key string) (bool, *utils.StringMapSet, string) {
	return true, nil, ""
}

func (r *hostRunner) Name() string {
	return "host"
}

func (r *hostRunner) BuildsRoot() string {
	return path.Join(core.GetWorkDir(), "builds")
}

func (r *hostRunner) TestsRoot() string {
	return path.Join(core.GetWorkDir(), "tests")
}

func (r *hostRunner) PrepareBuild(build *loader.AstBuild, cmdRunner loader.CmdRunner) *errors.JVSRuntimeResult {
	_, buildName := loader.ParseBuildName(build.Name)
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
		fileListContent, err := loader.GetCurSimulator().GetFileList(filelistItem...)
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
	if err := utils.WriteNewFile(path.Join(buildDir, "compile.sh"), loader.GetCurSimulator().CompileCmd()+" "+build.CompileOption()+filelistCompileOption); err != nil {
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

func (r *hostRunner) Build(build *loader.AstBuild, cmdRunner loader.CmdRunner) *errors.JVSRuntimeResult {
	_, buildName := loader.ParseBuildName(build.Name)
	buildDir := path.Join(r.BuildsRoot(), buildName)
	//create log file
	logFile, err := os.Create(path.Join(buildDir, buildName+".log"))
	defer logFile.Close()
	if err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	attr := loader.CmdAttr{WriteClosers: []io.WriteCloser{logFile},
		SetAttr: func(cmd *exec.Cmd) error {
			cmd.Dir = path.Join(r.BuildsRoot(), buildName)
			return nil
		}}
	res := cmdRunner(&attr, "bash", "run_compile.sh")
	return errors.NewJVSRuntimeResult(res.Status, res.GetMsg()+"\n", "path:"+buildDir)
}

func (r *hostRunner) PrepareTest(testCase *loader.AstTestCase, cmdRunner loader.CmdRunner) *errors.JVSRuntimeResult {
	_, buildName, testName, seed, groupsName := loader.ParseTestName(testCase.Name)
	testDir := path.Join(r.TestsRoot(), path.Join(groupsName...), buildName+"__"+testName, seed)
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
	if err := utils.WriteNewFile(path.Join(testDir, "sim.sh"), path.Join(buildName, loader.GetCurSimulator().SimCmd())+" "+testCase.SimOption()+" "+testCase.GetBuild().GetTestDiscoverer().TestCmd()+testName); err != nil {
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

func (r *hostRunner) RunTest(testCase *loader.AstTestCase, cmdRunner loader.CmdRunner) *errors.JVSRuntimeResult {
	_, buildName, testName, seed, groupsName := loader.ParseTestName(testCase.Name)
	testDir := path.Join(r.TestsRoot(), path.Join(groupsName...), buildName+"__"+testName, seed)
	//create log file
	logFile, err := os.Create(path.Join(testDir, buildName+"__"+testName+"__"+seed+".log"))
	defer logFile.Close()
	if err != nil {
		return errors.JVSRuntimeResultFail(err.Error())
	}
	attr := loader.CmdAttr{WriteClosers: []io.WriteCloser{logFile},
		SetAttr: func(cmd *exec.Cmd) error {
			cmd.Dir = testDir
			return nil
		}}
	res := cmdRunner(&attr, "bash", "run_sim.sh")
	return errors.NewJVSRuntimeResult(res.Status, res.GetMsg()+"\n", "path:"+testDir)
}

func init() {
	loader.RegisterRunner(newHostRunner)
}
