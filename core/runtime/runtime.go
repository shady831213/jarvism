package runtime

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/shady831213/jarvism/core/ast"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/options"
	"github.com/shady831213/jarvism/core/utils"
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

func hashFunc(s string) string {
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
	runTimeUnique = false
}

type runFlow struct {
	build     *ast.AstBuild
	testCases map[string]*ast.AstTestCase
	testWg    sync.WaitGroup
	cmdStdout *io.Writer
	buildDone chan *errors.JVSRuntimeResult
	testDone  chan *errors.JVSRuntimeResult
	ctx       context.Context
}

func newRunFlow(build *ast.AstBuild, cmdStdout *io.Writer, buildDone chan *errors.JVSRuntimeResult, testDone chan *errors.JVSRuntimeResult, ctx context.Context) *runFlow {
	inst := new(runFlow)
	inst.build = build
	inst.testWg = sync.WaitGroup{}
	inst.cmdStdout = cmdStdout
	inst.testCases = make(map[string]*ast.AstTestCase)
	inst.buildDone = buildDone
	inst.testDone = testDone
	inst.ctx = ctx
	return inst
}

type phase func() *errors.JVSRuntimeResult

func preparePhase(phaseName string, p phase) *errors.JVSRuntimeResult {
	PrintStatus(phaseName, utils.Blue("BEGIN"))
	result := p()
	if result == nil {
		result = errors.JVSRuntimeResultUnknown("No Result!")
		PrintStatus(phaseName, result.Error())
		return result
	}
	if result.Status != errors.JVSRuntimePass {
		PrintStatus(phaseName, result.Error())
	}
	return result
}

func runPhase(phaseName string, p phase) *errors.JVSRuntimeResult {
	result := p()
	if result == nil {
		result = errors.JVSRuntimeResultUnknown("No Result!")
		PrintStatus(phaseName, result.Error())
		return result
	}
	PrintStatus(phaseName, result.Error())
	return result
}

func (f *runFlow) cmdRunner(checkerPipeWriter io.WriteCloser) ast.CmdRunner {
	return func(attr *ast.CmdAttr, name string, arg ...string) (res *errors.JVSRuntimeResult) {
		cmd := exec.CommandContext(f.ctx, name, arg...)
		closers := make([]io.Closer, 0)
		defer func() {
			for _, c := range closers {
				if e := c.Close(); e != nil {
					res = errors.JVSRuntimeResultUnknown(e.Error())
				}
			}
		}()
		//set stdout
		writers := make([]io.Writer, 0)
		if *f.cmdStdout != nil {
			writers = append(writers, *f.cmdStdout)
		}
		//checker
		if checkerPipeWriter != nil {
			writers = append(writers, checkerPipeWriter)
			closers = append(closers, checkerPipeWriter)
		}
		//writeclosers in attr
		if attr != nil && attr.WriteClosers != nil {
			for _, wc := range attr.WriteClosers {
				writers = append(writers, wc)
				closers = append(closers, wc)
			}

		}

		fileAndStdoutWriter := io.MultiWriter(writers...)
		cmd.Stdout = fileAndStdoutWriter
		//set other attr
		if attr != nil && attr.SetAttr != nil {
			if err := attr.SetAttr(cmd); err != nil {
				return errors.JVSRuntimeResultUnknown(err.Error())
			}
		}
		if err := cmd.Run(); err != nil {
			return errors.JVSRuntimeResultUnknown(err.Error())
		}
		return errors.JVSRuntimeResultPass("")
	}
}

func (f *runFlow) prepareBuildPhase(build *ast.AstBuild) *errors.JVSRuntimeResult {
	return preparePhase(build.Name, func() *errors.JVSRuntimeResult {
		return ast.GetRunner().PrepareBuild(build, f.cmdRunner(nil))
	})
}

func (f *runFlow) checkPhase(checker ast.Checker) (*io.PipeWriter, func(), chan *errors.JVSRuntimeResult) {
	rd, wr := io.Pipe()
	checker.Input(rd)
	done := make(chan *errors.JVSRuntimeResult)
	goroutine := func() {
		defer close(done)
		select {
		case <-f.ctx.Done():
			done <- errors.JVSRuntimeResultUnknown("context canceled!")
		case done <- checker.Check():
			return
		}
	}
	return wr, goroutine, done
}

func (f *runFlow) buildPhase(build *ast.AstBuild) *errors.JVSRuntimeResult {
	return runPhase(build.Name, func() *errors.JVSRuntimeResult {
		wr, check, done := f.checkPhase(build.GetChecker())
		go check()
		status := errors.JVSRuntimePass
		execRes := ast.GetRunner().Build(build, f.cmdRunner(wr))
		if execRes.Status > status {
			status = execRes.Status
		}
		checkRes := <-done
		if checkRes.Status > status {
			status = checkRes.Status
		}
		return errors.NewJVSRuntimeResult(status, checkRes.GetMsg()+"\n", execRes.GetMsg())
	})
}

func (f *runFlow) prepareTestPhase(testCase *ast.AstTestCase) *errors.JVSRuntimeResult {
	return preparePhase(testCase.Name, func() *errors.JVSRuntimeResult {
		return ast.GetRunner().PrepareTest(testCase, f.cmdRunner(nil))
	})
}

func (f *runFlow) runTestPhase(testCase *ast.AstTestCase) *errors.JVSRuntimeResult {
	return runPhase(testCase.Name, func() *errors.JVSRuntimeResult {
		wr, check, done := f.checkPhase(testCase.GetChecker())
		go check()
		status := errors.JVSRuntimePass
		execRes := ast.GetRunner().RunTest(testCase, f.cmdRunner(wr))
		if execRes.Status > status {
			status = execRes.Status
		}
		checkRes := <-done
		if checkRes.Status > status {
			status = checkRes.Status
		}
		return errors.NewJVSRuntimeResult(status, checkRes.GetMsg()+"\n", execRes.GetMsg())
	})
}

func (f *runFlow) AddTest(test *ast.AstTestCase) int {
	test.Name = f.build.Name + "__" + test.Name
	test.SetBuild(f.build)
	if _, ok := f.testCases[test.Name]; !ok {
		f.testCases[test.Name] = test
		return 1
	}
	return 0
}

func (f *runFlow) run() {
	//run compile
	if !runTimeSimOnly {
		result := f.prepareBuildPhase(f.build)
		if result.Status != errors.JVSRuntimePass {
			f.buildDone <- result
			runTimeLimiter.get()
			return
		}
		result = f.buildPhase(f.build)
		if result.Status != errors.JVSRuntimePass {
			f.buildDone <- result
			runTimeLimiter.get()
			return
		}
		f.buildDone <- result
	}
	runTimeLimiter.get()

	//run tests
	for _, test := range f.testCases {
		f.testWg.Add(1)
		runTimeLimiter.put()
		go func(testCase *ast.AstTestCase) {
			defer f.testWg.Add(-1)
			defer runTimeLimiter.get()
			result := f.prepareTestPhase(testCase)
			if result.Status != errors.JVSRuntimePass {
				f.testDone <- result
				return
			}
			result = f.runTestPhase(testCase)
			f.testDone <- result
		}(test)
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
	buildDone                   chan *errors.JVSRuntimeResult
	testDone                    chan *errors.JVSRuntimeResult
	ctx                         context.Context
	cancel                      func()
}

func newRunTime(name string, group *ast.AstGroup) *runTime {
	r := new(runTime)
	r.Name = name
	r.runFlow = make(map[string]*runFlow)
	r.runtimeId = strings.Replace(time.Now().Format("20060102_150405.0000"), ".", "", 1)
	r.flowWg = sync.WaitGroup{}
	r.processingDone = make(chan bool)
	r.monitorDone = make(chan bool)
	r.buildDone = make(chan *errors.JVSRuntimeResult, 100)
	r.testDone = make(chan *errors.JVSRuntimeResult, 100)
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

func (r *runTime) createFlow(build *ast.AstBuild) *runFlow {
	var hash string
	if runTimeUnique {
		hash = hashFunc(r.runtimeId + build.GetRawSign())
	} else {
		hash = hashFunc(build.GetRawSign())
	}
	if _, ok := r.runFlow[hash]; !ok {
		newBuild := build.Clone()
		newBuild.Name = r.runtimeId + "__" + build.Name + "_" + hash
		r.runFlow[hash] = newRunFlow(newBuild, &r.cmdStdout, r.buildDone, r.testDone, r.ctx)
	}

	return r.runFlow[hash]
}

func (r *runTime) initSubTest(test *ast.AstTestCase) int {
	test.ParseArgs()
	flow := r.createFlow(test.GetBuild())
	cnt := 0
	for _, t := range test.GetTestCases() {
		cnt += flow.AddTest(t)
	}
	return cnt
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
			if a, _ := ast.GetJvsAstOption(arg); a != nil {
				_args = append(_args, arg)
			}
		}
	}
	return _args
}

func run(name string, cfg map[interface{}]interface{}, sc chan os.Signal) error {
	group := ast.NewAstGroup("Jarvis")
	if err := group.Parse(cfg); err != nil {
		return err
	}
	if err := group.Link(); err != nil {
		return err
	}
	r := newRunTime(name, group)
	logFile, err := setLog(r.runtimeId + ".log")
	defer func() {
		Println("logFile:" + logFile.Name())
		logFile.Close()
	}()
	if err != nil {
		return err
	}
	r.daemon(sc)
	return nil
}

func RunGroup(groupName string, args []string, sc chan os.Signal) error {
	return run(groupName, map[interface{}]interface{}{"args": filterAstArgs(args), "groups": []interface{}{groupName}}, sc)
}

func RunTest(testName, buildName string, args []string, sc chan os.Signal) error {
	return run(testName, map[interface{}]interface{}{"build": buildName,
		"args":  filterAstArgs(args),
		"tests": []interface{}{map[interface{}]interface{}{testName: nil}}}, sc)
}

func RunOnlyBuild(buildName string, args []string, sc chan os.Signal) error {
	return run(buildName, map[interface{}]interface{}{"build": buildName,
		"args": filterAstArgs(args)}, sc)
}

var runTimeMaxJob int
var runTimeSimOnly bool
var runTimeUnique bool

func init() {
	options.GetJvsOptions().IntVar(&runTimeMaxJob, "max_job", -1, "limit of runtime coroutines, default is unlimited.")
	options.GetJvsOptions().BoolVar(&runTimeSimOnly, "sim_only", false, "bypass compile and only run simulation, default is false.")
	options.GetJvsOptions().BoolVar(&runTimeUnique, "unique", false, "if set jobId(timestamp) will be included in hash, then builds and testcases will have unique name and be in unique dir.default is false.")

}
