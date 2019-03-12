package core

import (
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/shady831213/jarvism/core/options"
	"github.com/shady831213/jarvism/utils"
	"io"
	"math"
	"math/big"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
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
	buildDone chan *JVSTestResult
	testDone  chan *JVSTestResult
	ctx       context.Context
}

func newRunFlow(build *AstBuild, cmdStdout *io.Writer, buildDone chan *JVSTestResult, testDone chan *JVSTestResult, ctx context.Context) *runFlow {
	inst := new(runFlow)
	inst.build = build
	inst.testWg = sync.WaitGroup{}
	inst.cmdStdout = cmdStdout
	inst.List.Init()
	inst.buildDone = buildDone
	inst.testDone = testDone
	inst.ctx = ctx
	return inst
}

func (f *runFlow) cmdRunner(name string, arg ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdout = *f.cmdStdout
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		<-f.ctx.Done()
		cancel()
	}()
	return cmd.Wait()
}

func (f *runFlow) preparePhase(phaseName string, phase func() *JVSTestResult) *JVSTestResult {
	PrintStatus(phaseName, utils.Blue("BEGIN"))
	result := phase()
	if result == nil {
		result = JVSTestResultUnknown("No Result!")
		PrintStatus(phaseName, result.Error())
		return result
	}
	if result.status != JVSTestPass {
		PrintStatus(phaseName, result.Error())
	}
	return result
}

func (f *runFlow) runPhase(phaseName string, phase func() *JVSTestResult) *JVSTestResult {
	result := phase()
	if result == nil {
		result = JVSTestResultUnknown("No Result!")
		PrintStatus(phaseName, result.Error())
		return result
	}
	PrintStatus(phaseName, result.Error())
	return result
}

func (f *runFlow) prepareBuildPhase(build *AstBuild) *JVSTestResult {
	return f.preparePhase(build.Name, func() *JVSTestResult {
		return GetRunner().PrepareBuild(build, f.cmdRunner)
	})
}

func (f *runFlow) buildPhase(build *AstBuild) *JVSTestResult {
	return f.runPhase(build.Name, func() *JVSTestResult {
		return GetRunner().Build(build, f.cmdRunner)
	})
}

func (f *runFlow) prepareTestPhase(testCase *AstTestCase) *JVSTestResult {
	return f.preparePhase(testCase.Name, func() *JVSTestResult {
		return GetRunner().PrepareTest(testCase, f.cmdRunner)
	})
}

func (f *runFlow) runTestPhase(testCase *AstTestCase) *JVSTestResult {
	return f.runPhase(testCase.Name, func() *JVSTestResult {
		return GetRunner().RunTest(testCase, f.cmdRunner)
	})
}

func (f *runFlow) AddTest(test *AstTestCase) {
	test.Name = f.build.Name + "__" + test.Name
	test.Build = f.build
	f.PushBack(test)
}

func (f *runFlow) run() {
	//run compile
	if !runTimeSimOnly {
		result := f.prepareBuildPhase(f.build)
		if result.status != JVSTestPass {
			f.buildDone <- result
			runTimeLimiter.get()
			return
		}
		result = f.buildPhase(f.build)
		if result.status != JVSTestPass {
			f.buildDone <- result
			runTimeLimiter.get()
			return
		}
		f.buildDone <- result
	}
	runTimeLimiter.get()

	//run tests
	for e := f.Front(); e != nil; e = e.Next() {
		f.testWg.Add(1)
		runTimeLimiter.put()
		go func(testCase *AstTestCase) {
			defer f.testWg.Add(-1)
			defer runTimeLimiter.get()
			result := f.prepareTestPhase(testCase)
			if result.status != JVSTestPass {
				f.testDone <- result
				return
			}
			result = f.runTestPhase(testCase)
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
	buildDone                   chan *JVSTestResult
	testDone                    chan *JVSTestResult
	ctx                         context.Context
	cancel                      func()
}

func newRunTime(name string, group *astGroup) *runTime {
	r := new(runTime)
	r.Name = name
	r.runFlow = make(map[string]*runFlow)
	r.runtimeId = strings.Replace(time.Now().Format("20060102_150405.0000"), ".", "", 1)
	r.flowWg = sync.WaitGroup{}
	r.processingDone = make(chan bool)
	r.monitorDone = make(chan bool)
	r.buildDone = make(chan *JVSTestResult, 100)
	r.testDone = make(chan *JVSTestResult, 100)
	ctx := context.Background()
	r.ctx, r.cancel = context.WithCancel(ctx)
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
		r.cmdStdout = &stdout{}
	}
	return r
}

func (r *runTime) createFlow(build *AstBuild) *runFlow {
	hash := hash(build.Name + build.compileItems.preAction + build.compileItems.option.GetString() + build.compileItems.postAction)
	if _, ok := r.runFlow[hash]; !ok {
		newBuild := build.Clone()
		newBuild.Name = r.runtimeId + "__" + build.Name + "_" + hash
		r.runFlow[hash] = newRunFlow(newBuild, &r.cmdStdout, r.buildDone, r.testDone, r.ctx)
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
	defer func() {
		close(r.buildDone)
		close(r.testDone)
	}()
	for _, f := range r.runFlow {
		r.flowWg.Add(1)
		runTimeLimiter.put()
		go func(flow *runFlow) {
			defer r.flowWg.Add(-1)
			flow.run()
		}(f)
	}
	r.flowWg.Wait()
	r.cancel()
}

func (r *runTime) exit() {
	r.monitorDone <- true
	close(r.monitorDone)
	r.processingDone <- true
	close(r.processingDone)
	runTimeFinish()
	Println("exit!")
}

func (r *runTime) signalHandler(sc chan os.Signal) {
	if sc != nil {
		signal.Notify(sc, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case s := <-sc:
			Println("receive signal" + s.String())
			r.cancel()
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *runTime) daemon(sc chan os.Signal) {

	defer r.exit()
	var status string

	// run

	//monitor status
	go PrintProccessing(utils.Brown)("Jarvism is running", &status, r.processingDone)
	go statusMonitor(&status, len(r.runFlow), r.totalTest, r.buildDone, r.testDone, r.monitorDone)

	//monitor signals and run
	go r.signalHandler(sc)
	r.run()
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

func run(name string, cfg map[interface{}]interface{}, sc chan os.Signal) error {
	group := newAstGroup("Jarvis")
	if err := group.Parse(cfg); err != nil {
		return err
	}
	if err := group.Link(); err != nil {
		return err
	}
	newRunTime(name, group).daemon(sc)
	return nil
}

func RunGroup(group *astGroup, args []string, sc chan os.Signal) error {
	return run(group.Name, map[interface{}]interface{}{"args": filterAstArgs(args), "groups": []interface{}{group.Name}}, sc)
}

func RunTest(testName, buildName string, args []string, sc chan os.Signal) error {
	return run(testName, map[interface{}]interface{}{"build": buildName,
		"args":  filterAstArgs(args),
		"tests": map[interface{}]interface{}{testName: nil}}, sc)
}

func RunOnlyBuild(buildName string, args []string, sc chan os.Signal) error {
	return run(buildName, map[interface{}]interface{}{"build": buildName,
		"args": filterAstArgs(args)}, sc)
}

var runTimeMaxJob int
var runTimeSimOnly bool

func init() {
	options.GetJvsOptions().IntVar(&runTimeMaxJob, "max_job", -1, "limit of runtime coroutines, default is unlimited.")
	options.GetJvsOptions().BoolVar(&runTimeSimOnly, "sim_only", false, "bypass compile and only run simulation, default is false.")
}
