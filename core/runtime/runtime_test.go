package runtime_test

import (
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/runtime"
	_ "github.com/shady831213/jarvism/simulators"
	_ "github.com/shady831213/jarvism/testDiscoverers"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"
)

type testRunner struct {
}

func (r *testRunner) Name() string {
	return "test"
}

func (r *testRunner) PrepareBuild(build *ast.AstBuild, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", " "); err != nil {
		return errors.JVSTestResultFail(err.Error())
	}
	return errors.JVSTestResultPass("")
}

func (r *testRunner) Build(build *ast.AstBuild, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", " build build ", build.Name); err != nil {
		return errors.JVSTestResultFail(err.Error())
	}
	return errors.JVSTestResultPass("")
}

func (r *testRunner) PrepareTest(testCase *ast.AstTestCase, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", ""); err != nil {
		return errors.JVSTestResultFail(err.Error())
	}
	return errors.JVSTestResultPass("")
}

func (r *testRunner) RunTest(testCase *ast.AstTestCase, cmdRunner func(string, ...string) error) *errors.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", " run test ", testCase.Name); err != nil {
		return errors.JVSTestResultFail(err.Error())
	}
	return errors.JVSTestResultPass("")
}

func TestGroup(t *testing.T) {
	if err := runtime.RunGroup(ast.GetJvsAstRoot().GetGroup("group1"), []string{"-sim_only"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := runtime.RunGroup(ast.GetJvsAstRoot().GetGroup("group2"), []string{"-max_job " + strconv.Itoa(rand.Intn(2)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := runtime.RunGroup(ast.GetJvsAstRoot().GetGroup("group3"), []string{"-max_job " + strconv.Itoa(rand.Intn(50)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestSingleTest(t *testing.T) {
	if err := runtime.RunTest("test1", "build1", []string{"-seed 1"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestSingleRepeatTest(t *testing.T) {
	if err := runtime.RunTest("test1", "build1", []string{"-repeat 10", "-max_job " + strconv.Itoa(rand.Intn(15)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestRunOnlyBuild(t *testing.T) {
	if err := runtime.RunOnlyBuild("build1", []string{"-test_phase jarvis"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestInterrupt(t *testing.T) {
	sc := make(chan os.Signal)
	defer close(sc)
	go func() {
		var stopChan = make(chan os.Signal)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT)
		sc <- <-stopChan
	}()
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Int63n(400)) * time.Millisecond)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()

		if err := runtime.RunTest("test1", "build1", []string{"-repeat 10", "-max_job " + strconv.Itoa(rand.Intn(15)+1)}, sc); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

func init() {
	ast.RegisterRunner(new(testRunner))
	os.Setenv("JVS_PRJ_HOME", "testFiles")
	cfg, err := ast.Lex("testFiles/build.yaml")
	if err != nil {
		panic(err)
	}
	err = ast.Parse(cfg)
	if err != nil {
		panic(err)
	}
}
