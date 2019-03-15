package runtime_test

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/runtime"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func setup() {
	os.Symlink(path.Join(core.CorePath(), "runtime", "testFiles", "jarvism_plugins"), "/tmp/jarvism_plugins")
}

func tearDonw() {
	os.RemoveAll(path.Join(loader.GetWorkDir(), "JarvismLog"))
	os.RemoveAll("/tmp/jarvism_plugins")
}

func TestGroup(t *testing.T) {
	setup()
	if err := runtime.RunGroup("group1", []string{"-sim_only"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := runtime.RunGroup("group2", []string{"-max_job " + strconv.Itoa(rand.Intn(2)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := runtime.RunGroup("group3", []string{"-max_job " + strconv.Itoa(rand.Intn(50)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	tearDonw()
}

func TestSingleTest(t *testing.T) {
	setup()
	if err := runtime.RunTest("test1", "build1", []string{"-seed 1"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	tearDonw()

}

func TestSingleRepeatTest(t *testing.T) {
	setup()
	if err := runtime.RunTest("test1", "build1", []string{"-repeat 10", "-max_job " + strconv.Itoa(rand.Intn(15)+1)}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	tearDonw()

}

func TestRunOnlyBuild(t *testing.T) {
	setup()
	if err := runtime.RunOnlyBuild("build1", []string{"-test_phase jarvis"}, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}
	tearDonw()

}

func TestInterrupt(t *testing.T) {
	setup()
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
	tearDonw()
}

func init() {
	//build and load plugin
	abs, _ := filepath.Abs("testFiles")
	os.Setenv("JVS_PRJ_HOME", abs)
	os.Setenv("JVS_PLUGINS_HOME", "/tmp/jarvism_plugins")
	os.Symlink(path.Join(core.CorePath(), "runtime", "testFiles", "jarvism_plugins"), "/tmp/jarvism_plugins")
	cfg, err := loader.Lex(path.Join(abs, "build.yaml"))
	if err != nil {
		panic(err)
	}
	err = loader.Parse(cfg)
	if err != nil {
		panic(err)
	}
}
