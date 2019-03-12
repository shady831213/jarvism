package runtime

import (
	"github.com/shady831213/jarvism/core/ast"
	"strconv"
	"testing"
)

func setUp(name string, cfg map[interface{}]interface{}) (*runTime, error) {
	group := ast.NewAstGroup("Jarvis")
	if err := group.Parse(cfg); err != nil {
		return nil, err
	}
	if err := group.Link(); err != nil {
		return nil, err
	}
	return newRunTime(name, group), nil
}

func setUpGroup(group *ast.AstGroup, args []string) (*runTime, error) {
	return setUp(group.Name, map[interface{}]interface{}{"args": filterAstArgs(args), "groups": []interface{}{group.Name}})
}

func setUpTest(testName, buildName string, args []string) (*runTime, error) {
	return setUp(testName, map[interface{}]interface{}{"build": buildName,
		"args":  filterAstArgs(args),
		"tests": map[interface{}]interface{}{testName: nil}})
}

func setUpOnlyBuild(buildName string, args []string) (*runTime, error) {
	return setUp(buildName, map[interface{}]interface{}{"build": buildName,
		"args": filterAstArgs(args)})
}

func TestGroupSetup(t *testing.T) {
	r, err := setUpGroup(ast.GetJvsAstRoot().GetGroup("group1"), []string{"-max_job 5"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(r.runFlow) != 1 {
		t.Error("expect 1 runFlow but get " + string(len(r.runFlow)))
		t.FailNow()
	}
	if runTimeMaxJob != 5 {
		t.Error("runTimeMaxJob expect 5, but get " + strconv.Itoa(runTimeMaxJob))
		t.FailNow()
	}
	runTimeFinish()

	r, err = setUpGroup(ast.GetJvsAstRoot().GetGroup("group2"), []string{"-sim_only"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(r.runFlow) != 2 {
		t.Error("expect 2 runFlow but get " + string(len(r.runFlow)))
		t.FailNow()
	}
	if runTimeSimOnly != true {
		t.Error("runTimeSimOnly expect true, but get false")
		t.FailNow()
	}
	runTimeFinish()

	r, err = setUpGroup(ast.GetJvsAstRoot().GetGroup("group3"), []string{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(r.runFlow) != 2 {
		t.Error("expect 2 runFlow but get " + string(len(r.runFlow)))
		t.FailNow()
	}
	runTimeFinish()
}

func TestSingleTestSetup(t *testing.T) {
	defer runTimeFinish()
	r, err := setUpTest("test1", "build1", []string{"-seed 1"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if r.cmdStdout == nil {
		t.Error("when running single test, expect stdout is open but closed")
		t.FailNow()
	}
	if runTimeMaxJob != -1 {
		t.Error("runTimeMaxJob expect -1, but get " + strconv.Itoa(runTimeMaxJob))
		t.FailNow()
	}
}

func TestRunOnlyBuildSetup(t *testing.T) {
	defer runTimeFinish()
	r, err := setUpOnlyBuild("build1", []string{"-test_phase jarvis", "-max_job 10"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if r.cmdStdout == nil {
		t.Error("when running only build, expect stdout is open but closed")
		t.FailNow()
	}
	if runTimeMaxJob != 10 {
		t.Error("runTimeMaxJob expect 10, but get " + strconv.Itoa(runTimeMaxJob))
		t.FailNow()
	}

}

