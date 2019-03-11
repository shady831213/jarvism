package core_test

import (
	"github.com/shady831213/jarvism/core"
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
	return "testRunner"
}

func (r *testRunner) PrepareBuild(build *core.AstBuild, cmdRunner func(string, ...string) error) *core.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", " "); err != nil {
		return core.JVSTestResultFail(err.Error())
	}
	return core.JVSTestResultPass("")
}

func (r *testRunner) Build(build *core.AstBuild, cmdRunner func(string, ...string) error) *core.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", " build build ", build.Name); err != nil {
		return core.JVSTestResultFail(err.Error())
	}
	return core.JVSTestResultPass("")
}

func (r *testRunner) PrepareTest(testCase *core.AstTestCase, cmdRunner func(string, ...string) error) *core.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", ""); err != nil {
		return core.JVSTestResultFail(err.Error())
	}
	return core.JVSTestResultPass("")
}

func (r *testRunner) RunTest(testCase *core.AstTestCase, cmdRunner func(string, ...string) error) *core.JVSTestResult {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
	if err := cmdRunner("echo", " run test ", testCase.Name); err != nil {
		return core.JVSTestResultFail(err.Error())
	}
	return core.JVSTestResultPass("")
}

func TestGroup(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunGroup(core.GetJvsAstRoot().GetGroup("group1"), []string{"-sim_only"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := core.RunGroup(core.GetJvsAstRoot().GetGroup("group2"), []string{"-max_job " + strconv.Itoa(rand.Intn(2)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := core.RunGroup(core.GetJvsAstRoot().GetGroup("group3"), []string{"-max_job " + strconv.Itoa(rand.Intn(50)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestSingleTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunTest("test1", "build1", []string{"-seed 1"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestSingleRepeatTest(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunTest("test1", "build1", []string{"-repeat 10", "-max_job " + strconv.Itoa(rand.Intn(15)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestRunOnlyBuild(t *testing.T) {
	core.SetRunner(new(testRunner))
	if err := core.RunOnlyBuild("build1", []string{"-test_phase jarvis"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestInterrupt(t *testing.T) {
	core.SetRunner(new(testRunner))
	sc := make(chan os.Signal)
	defer close(sc)
	go func() {
		var stopChan = make(chan os.Signal)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT)
		sc <- <-stopChan
	}()
	for i := 0; i < 10; i ++ {
		go func() {
			time.Sleep(time.Duration(rand.Int63n(400)) * time.Millisecond)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()

		if err := core.RunTest("test1", "build1", []string{"-repeat 10", "-max_job " + strconv.Itoa(rand.Intn(15)+1)}, sc); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
