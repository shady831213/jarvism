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

type runTimeJobLimiter struct {
	maxJob chan bool
}

func (l *runTimeJobLimiter) put() {
	if l.maxJob != nil {
		l.maxJob <- true
	}
}

func (l *runTimeJobLimiter) get() {
	if l.maxJob != nil {
		<-l.maxJob
	}
}

func (l *runTimeJobLimiter) close() {
	if l.maxJob != nil {
		close(l.maxJob)
	}
}

var runTimeLimiter runTimeJobLimiter

func runTimeFinish() {
	runTimeLimiter.close()
	runTimeMaxJob = -1
	runTimeSimOnly = false
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
	PrintStatus(build.Name, utils.Blue("BEGIN"))
	if err := GetRunner().PrepareBuild(build, cmdStdout); err != nil {
		PrintStatus(build.Name, utils.Red("FAIL"))
		return err
	}
	return nil
}

func (f *runFlow) buildPhase(build *AstBuild, cmdStdout *io.Writer) error {
	if err := GetRunner().Build(build, cmdStdout); err != nil {
		PrintStatus(build.Name, utils.Red("FAIL"))
		return err
	}
	PrintStatus(build.Name, utils.Green("DONE"))
	return nil
}

func (f *runFlow) prepareTestPhase(testCase *AstTestCase, cmdStdout *io.Writer) *JVSTestResult {
	PrintStatus(testCase.Name, utils.Blue("BEGIN"))
	result := GetRunner().PrepareTest(testCase, cmdStdout)
	if result == nil {
		result = JVSTestResultUnknown("No Result!")
		PrintStatus(testCase.Name, result.Error())
		return result
	}
	return result
}

func (f *runFlow) runTestPhase(testCase *AstTestCase, cmdStdout *io.Writer) *JVSTestResult {
	result := GetRunner().RunTest(testCase, cmdStdout)
	if result == nil {
		result = JVSTestResultUnknown("No Result!")
		PrintStatus(testCase.Name, result.Error())
		return result
	}
	PrintStatus(testCase.Name, result.Error())
	return result
}

func (f *runFlow) AddTest(test *AstTestCase) {
	test.Name = f.build.Name + "__" + test.Name
	test.Build = f.build
	f.PushBack(test)
}

func (f *runFlow) run() {
	//run compile
	if !runTimeSimOnly {
		if err := f.prepareBuildPhase(f.build, f.cmdStdout); err != nil {
			fmt.Println(utils.Red(err.Error()))
			f.buildDone <- err
			runTimeLimiter.get()
			return
		}
		if err := f.buildPhase(f.build, f.cmdStdout); err != nil {
			fmt.Println(utils.Red(err.Error()))
			f.buildDone <- err
			runTimeLimiter.get()
			return
		}
		f.buildDone <- nil
	}
	runTimeLimiter.get()

	//run tests
	for e := f.Front(); e != nil; e = e.Next() {
		runTimeLimiter.put()
		f.testWg.Add(1)
		go func(testCase *AstTestCase) {
			defer f.testWg.Add(-1)
			defer runTimeLimiter.get()
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
	if runTimeMaxJob > 0 {
		runTimeLimiter = runTimeJobLimiter{make(chan bool, runTimeMaxJob)}
	} else {
		runTimeLimiter = runTimeJobLimiter{nil}
	}

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
	hash := hash(build.Name + build.compileItems.preAction + build.compileItems.option.GetString() + build.compileItems.postAction)
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
	defer runTimeFinish()
	var status string
	go PrintProccessing(utils.Brown)("Jarvism is running", &status, r.processingDone)
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
		runTimeLimiter.put()
		r.flowWg.Add(1)
		go func(flow *runFlow) {
			defer r.flowWg.Add(-1)
			flow.run()
		}(f)
	}
	r.flowWg.Wait()
}

func filterAstArgs(args []string) []interface{} {
	_args := make([]interface{}, 0)
	if args != nil {
		for _, arg := range args {
			//Parse all args and only pass the jvsAstOption to Ast
			if a, _ := getJvsAstOption(arg); a != nil {
				_args = append(_args, arg)
			}
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
	return run(group.Name, map[interface{}]interface{}{"args": filterAstArgs(args), "groups": []interface{}{group.Name}})
}

func RunTest(testName, buildName string, args []string) error {
	return run(testName, map[interface{}]interface{}{"build": buildName,
		"args":  filterAstArgs(args),
		"tests": map[interface{}]interface{}{testName: nil}})
}

func RunOnlyBuild(buildName string, args []string) error {
	return run(buildName, map[interface{}]interface{}{"build": buildName,
		"args": filterAstArgs(args)})
}
