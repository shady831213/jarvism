package core

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/shady831213/jarvism/utils"
	"io"
	"math"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"
)

func hash(s string) string {
	h := new(big.Int).SetBytes(sha256.New().Sum(([]byte)(s)))
	mb := big.NewInt(math.MaxInt64)
	h.Mod(h, mb)
	return hex.EncodeToString(h.Bytes())
}

type runTimeOpts interface {
	GetName() string
	//bottom-up search
	GetBuild() *AstBuild
	//top-down search
	GetTestCases() []*AstTestCase
	ParseArgs()
}

type runFlow struct {
	build *AstBuild
	list.List
	testWg    sync.WaitGroup
	cmdStdout *io.Writer
	buildDone chan error
	testDone  chan *JVSTestResult
}

func runBuildFlowPhase(phase func(build *AstBuild, cmdStdout *io.Writer) error, phaseName string) func(*AstBuild, *io.Writer) error {
	return func(build *AstBuild, cmdStdout *io.Writer) error {
		PrintStatus(phaseName+build.Name, utils.Blue("BEGIN"))
		if err := phase(build, cmdStdout); err != nil {
			PrintStatus(phaseName+build.Name, utils.Red("FAIL"))
			return err
		}
		PrintStatus(phaseName+build.Name, utils.Green("DONE"))
		return nil
	}
}

func runTestFlowPhase(phase func(testCase *AstTestCase, cmdStdout *io.Writer) *JVSTestResult, phaseName string) func(*AstTestCase, *io.Writer) *JVSTestResult {
	return func(testCase *AstTestCase, cmdStdout *io.Writer) *JVSTestResult {
		PrintStatus(phaseName+testCase.Name, utils.Blue("BEGIN"))
		result := phase(testCase, cmdStdout)
		if result == nil {
			result = JVSTestResultUnknown("No Result!")
			PrintStatus(phaseName+testCase.Name, result.Error())
			return result
		}
		PrintStatus(phaseName+testCase.Name, result.Error())
		return result
	}
}

func newRunFlow(build *AstBuild, cmdStdout *io.Writer, buildDone chan error, testDone chan *JVSTestResult) *runFlow {
	inst := new(runFlow)
	inst.build = build
	inst.testWg = sync.WaitGroup{}
	inst.cmdStdout = cmdStdout
	inst.List.Init()
	inst.buildDone = buildDone
	inst.testDone = testDone
	return inst
}

func (f *runFlow) prepareBuildPhase(build *AstBuild, cmdStdout *io.Writer) error {
	return runBuildFlowPhase(GetRunner().PrepareBuild, "Prepare Build ")(build, cmdStdout)
}

func (f *runFlow) buildPhase(build *AstBuild, cmdStdout *io.Writer) error {
	return runBuildFlowPhase(GetRunner().Build, "Build Build ")(build, cmdStdout)
}

func (f *runFlow) prepareTestPhase(testCase *AstTestCase, cmdStdout *io.Writer) *JVSTestResult {
	return runTestFlowPhase(GetRunner().PrepareTest, "Prepare Test ")(testCase, cmdStdout)
}

func (f *runFlow) runTestPhase(testCase *AstTestCase, cmdStdout *io.Writer) *JVSTestResult {
	return runTestFlowPhase(GetRunner().RunTest, "Run Test ")(testCase, cmdStdout)
}

func (f *runFlow) AddTest(test *AstTestCase) {
	test.Name = f.build.Name + "__" + test.Name
	test.Build = f.build
	f.PushBack(test)
}

func (f *runFlow) run() {
	if err := f.prepareBuildPhase(f.build, f.cmdStdout); err != nil {
		fmt.Println(utils.Red(err.Error()))
		f.buildDone <- err
		return
	}
	if err := f.buildPhase(f.build, f.cmdStdout); err != nil {
		fmt.Println(utils.Red(err.Error()))
		f.buildDone <- err
		return
	}
	f.buildDone <- nil
	for e := f.Front(); e != nil; e = e.Next() {
		f.testWg.Add(1)
		go func(testCase *AstTestCase) {
			defer f.testWg.Add(-1)
			result := f.prepareTestPhase(testCase, f.cmdStdout)
			if result.status != JVSTestPass {
				f.testDone <- result
				return
			}
			result = f.runTestPhase(testCase, f.cmdStdout)
			f.testDone <- result
		}(e.Value.(*AstTestCase))
	}
	f.testWg.Wait()
}

type runTime struct {
	cmdStdout                   io.Writer
	runtimeId                   string
	Name                        string
	totalTest                   int
	runFlow                     map[string]*runFlow
	flowWg                      sync.WaitGroup
	processingDone, monitorDone chan bool
	buildDone                   chan error
	testDone                    chan *JVSTestResult
}

func newRunTime(name string, group *astGroup) *runTime {
	r := new(runTime)
	r.Name = name
	r.runFlow = make(map[string]*runFlow)
	r.runtimeId = strings.Replace(time.Now().Format("20060102_150405.0000"), ".", "", 1)
	r.flowWg = sync.WaitGroup{}
	r.processingDone = make(chan bool)
	r.monitorDone = make(chan bool)
	r.buildDone = make(chan error, 100)
	r.testDone = make(chan *JVSTestResult, 100)

	testcases := group.GetTestCases()
	r.totalTest = 0
	for _, test := range testcases {
		r.totalTest += r.initSubTest(test)
	}
	//build only
	if r.totalTest == 0 {
		group.ParseArgs()
		r.createFlow(group.GetBuild())
	}
	if r.totalTest <= 1 {
		r.cmdStdout = os.Stdout
	}
	return r
}

func (r *runTime) createFlow(build *AstBuild) *runFlow {
	hash := hash(build.Name + build.compileItems.preAction + build.compileItems.option.GetString() + build.compileItems.postAction + r.runtimeId)
	if _, ok := r.runFlow[hash]; !ok {
		newBuild := build.Clone()
		newBuild.Name = r.runtimeId + "__" + build.Name + "_" + hash
		r.runFlow[hash] = newRunFlow(newBuild, &r.cmdStdout, r.buildDone, r.testDone)
	}

	return r.runFlow[hash]
}

func (r *runTime) initSubTest(test *AstTestCase) int {
	test.ParseArgs()
	flow := r.createFlow(test.GetBuild())
	testcases := test.GetTestCases()
	for _, t := range test.GetTestCases() {
		flow.AddTest(t)
	}
	return len(testcases)
}

func (r *runTime) run() {
	var status string
	go PrintProccessing(utils.Brown)(r.Name+" is running", &status, r.processingDone)
	defer func() {
		r.processingDone <- true
		close(r.processingDone)
	}()
	go statusMonitor(&status, len(r.runFlow), r.totalTest, r.buildDone, r.testDone, r.monitorDone)
	defer func() {
		r.monitorDone <- true
		close(r.monitorDone)
		close(r.buildDone)
		close(r.testDone)
	}()
	for _, f := range r.runFlow {
		r.flowWg.Add(1)
		go func(flow *runFlow) {
			defer r.flowWg.Add(-1)
			flow.run()
		}(f)
	}
	r.flowWg.Wait()
}

func convertArgs(args []string) []interface{} {
	_args := make([]interface{}, 0)
	if args != nil {
		for _, a := range args {
			_args = append(_args, a)
		}
	}
	return _args
}

func run(name string, cfg map[interface{}]interface{}) error {
	group := newAstGroup("Jarvis")
	if err := group.Parse(cfg); err != nil {
		return err
	}
	if err := group.Link(); err != nil {
		return err
	}
	newRunTime(name, group).run()
	return nil
}

func RunGroup(group *astGroup, args []string) error {
	return run(group.Name, map[interface{}]interface{}{"args": convertArgs(args), "groups": []interface{}{group.Name}})
}

func RunTest(testName, buildName string, args []string) error {
	return run(testName, map[interface{}]interface{}{"build": buildName,
		"args":  convertArgs(args),
		"tests": map[interface{}]interface{}{testName: nil}})
}

func RunOnlyBuild(buildName string, args []string) error {
	return run(buildName, map[interface{}]interface{}{"build": buildName,
		"args": convertArgs(args)})
}
