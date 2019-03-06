package core

import (
	"container/list"
	"fmt"
	"github.com/shady831213/jarvisSim/utils"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type runTimeOpts interface {
	GetName() string
	//bottom-up search
	GetBuild() *AstBuild
	//top-down search
	GetTestCases() []*AstTestCase
}

type runFlow struct {
	build *AstBuild
	list.List
	testWg    sync.WaitGroup
	cmdStdout io.Writer
}

func runBuildFlowPhase(phase func(build *AstBuild, cmdStdout io.Writer) error, phaseName string) func(*AstBuild, io.Writer) error {
	return func(build *AstBuild, cmdStdout io.Writer) error {
		utils.PrintStatus(utils.Blue, utils.Blue)(phaseName+build.Name, "BEGIN")
		if err := phase(build, cmdStdout); err != nil {
			utils.PrintStatus(utils.Blue, utils.Red)(phaseName+build.Name, "FAIL")
			return err
		}
		utils.PrintStatus(utils.Blue, utils.Green)(phaseName+build.Name, "DONE")
		return nil
	}
}

func runTestFlowPhase(phase func(testCase *AstTestCase, cmdStdout io.Writer) error, phaseName string) func(*AstTestCase, io.Writer) error {
	return func(testCase *AstTestCase, cmdStdout io.Writer) error {
		utils.PrintStatus(utils.Blue, utils.Blue)(phaseName+testCase.Name, "BEGIN")
		if err := phase(testCase, cmdStdout); err != nil {
			utils.PrintStatus(utils.Blue, utils.Red)(phaseName+testCase.Name, "FAIL")
			return err
		}
		utils.PrintStatus(utils.Blue, utils.Green)(phaseName+testCase.Name, "DONE")
		return nil
	}
}

func newRunFlow(build *AstBuild, cmdStdout io.Writer) *runFlow {
	inst := new(runFlow)
	inst.build = build
	inst.testWg = sync.WaitGroup{}
	inst.cmdStdout = cmdStdout
	inst.List.Init()
	return inst
}

func (f *runFlow) prepareBuildPhase(build *AstBuild, cmdStdout io.Writer) error {
	return runBuildFlowPhase(GetRunner().PrepareBuild, "Prepare Build ")(build, cmdStdout)
}

func (f *runFlow) buildPhase(build *AstBuild, cmdStdout io.Writer) error {
	return runBuildFlowPhase(GetRunner().Build, "Build Build ")(build, cmdStdout)
}

func (f *runFlow) prepareTestPhase(testCase *AstTestCase, cmdStdout io.Writer) error {
	return runTestFlowPhase(GetRunner().PrepareTest, "Prepare Test ")(testCase, cmdStdout)
}

func (f *runFlow) runTestPhase(testCase *AstTestCase, cmdStdout io.Writer) error {
	return runTestFlowPhase(GetRunner().RunTest, "Run Test ")(testCase, cmdStdout)
}

func (f *runFlow) run() {
	if err := f.prepareBuildPhase(f.build, f.cmdStdout); err != nil {
		fmt.Println(utils.Red(err.Error()))
		return
	}
	if err := f.buildPhase(f.build, f.cmdStdout); err != nil {
		fmt.Println(utils.Red(err.Error()))
		return
	}
	for e := f.Front(); e != nil; e = e.Next() {
		f.testWg.Add(1)
		go func(testCase *AstTestCase) {
			defer f.testWg.Add(-1)
			if err := f.prepareTestPhase(testCase, f.cmdStdout); err != nil {
				fmt.Println(utils.Red(err.Error()))
				return
			}
			if err := f.runTestPhase(testCase, f.cmdStdout); err != nil {
				fmt.Println(utils.Red(err.Error()))
			}
		}(e.Value.(*AstTestCase))
	}
	f.testWg.Wait()
}

type runTime struct {
	cmdStdout io.Writer
	timeStamp string
	Name      string
	runFlow   map[string]*runFlow
	flowWg    sync.WaitGroup
	done      chan bool
}

func (r *runTime) init(opts runTimeOpts) {
	r.Name = opts.GetName()
	r.runFlow = make(map[string]*runFlow)
	r.timeStamp = strings.Replace(time.Now().Format("20060102_150405.0000"), ".", "", 1)
	testcases := opts.GetTestCases()
	if len(testcases) <= 1 {
		r.cmdStdout = os.Stdout
	}
	r.flowWg = sync.WaitGroup{}
	r.done = make(chan bool)
	for _, test := range testcases {
		test.Name = r.timeStamp + "__" + test.Name
		if _, ok := r.runFlow[test.GetBuild().Name]; !ok {
			r.runFlow[test.GetBuild().Name] = newRunFlow(test.GetBuild(), r.cmdStdout)
		}
		r.runFlow[test.GetBuild().Name].PushBack(test)
	}
}

func (r *runTime) run() {
	go utils.PrintProccessing(utils.Blue)(r.Name+" is running", r.done)
	defer func() { r.done <- true }()
	for _, f := range r.runFlow {
		r.flowWg.Add(1)
		go func(flow *runFlow) {
			defer r.flowWg.Add(-1)
			f.run()
		}(f)
	}
	r.flowWg.Wait()
}

func Run(opts runTimeOpts) {
	r := new(runTime)
	r.init(opts)
	r.run()
}
