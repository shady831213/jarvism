package runtime_test

import (
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/runtime"
	_ "github.com/shady831213/jarvism/simulators"
	_ "github.com/shady831213/jarvism/testDiscoverers"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"plugin"
	"strconv"
	"syscall"
	"testing"
	"time"
)

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
	os.RemoveAll(path.Join(ast.GetWorkDir(), "JarvismLog"))
}

func TestSingleTest(t *testing.T) {
	if err := runtime.RunTest("test1", "build1", []string{"-seed 1"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	os.RemoveAll(path.Join(ast.GetWorkDir(), "JarvismLog"))

}

func TestSingleRepeatTest(t *testing.T) {
	if err := runtime.RunTest("test1", "build1", []string{"-repeat 10", "-max_job " + strconv.Itoa(rand.Intn(15)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	os.RemoveAll(path.Join(ast.GetWorkDir(), "JarvismLog"))

}

func TestRunOnlyBuild(t *testing.T) {
	if err := runtime.RunOnlyBuild("build1", []string{"-test_phase jarvis"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	os.RemoveAll(path.Join(ast.GetWorkDir(), "JarvismLog"))

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
	os.RemoveAll(path.Join(ast.GetWorkDir(), "JarvismLog"))
}

func init() {
	//build and load plugin
	cmd := exec.Command("go", "build", "-o", "testRunner.so", "-buildmode", "plugin", "./testRunner")
	cmd.Dir = "testFiles/plugins/runners"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("build plugin fail:", err)
		return
	}
	_, err := plugin.Open("testFiles/plugins/runners/testRunner.so")
	if err != nil {
		log.Println("open plugin err:", err, "testFiles/plugins/runners/testRunner.so")
		return
	}

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
