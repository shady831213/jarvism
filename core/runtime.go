package core

import (
	"container/list"
	"strconv"
	"sync"
	"time"
)

type runTimeOpts interface {
	GetName() string
	//bottom-up search
	GetBuild() *astBuild
	//top-down search
	GetTestCases() []*AstTestCase
}

type runFlow struct {
	build     *astBuild
	buildDone chan bool
	list.List
	testWg sync.WaitGroup
}

func newRunFlow(build *astBuild) *runFlow {
	inst := new(runFlow)
	inst.build = build
	inst.buildDone = make(chan bool, 1)
	inst.testWg = sync.WaitGroup{}
	inst.List.Init()
	return inst
}

type runTime struct {
	timeStamp string
	Name      string
	runFlow   map[string]*runFlow
}

func (r *runTime) Init(opts runTimeOpts) {
	r.Name = opts.GetName()
	r.timeStamp = strconv.FormatInt(time.Now().UnixNano(), 16)
	testcases := opts.GetTestCases()
	for _, test := range testcases {
		test.Name = r.timeStamp + "__" + test.Name
		if _, ok := r.runFlow[test.GetBuild().Name]; !ok {
			r.runFlow[test.GetBuild().Name] = newRunFlow(test.GetBuild())
		}
		r.runFlow[test.GetBuild().Name].PushBack(test)
	}
}
