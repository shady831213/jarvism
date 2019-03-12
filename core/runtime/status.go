package runtime

import (
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/utils"
	"strconv"
)

type statusCnt struct {
	cnts   map[errors.JVSTestStatus]int
	total  int
	keys   []errors.JVSTestStatus
	name   string
	finish bool
}

func newStatusCnt(name string, total int) *statusCnt {
	inst := new(statusCnt)
	inst.name = name
	inst.cnts = make(map[errors.JVSTestStatus]int)
	inst.keys = make([]errors.JVSTestStatus, 0)
	inst.total = total
	inst.cnts[errors.JVSTestPass] = 0
	inst.cnts[errors.JVSTestFail] = 0
	inst.cnts[errors.JVSTestWarning] = 0
	inst.cnts[errors.JVSTestUnknown] = 0
	inst.keys = append(inst.keys, errors.JVSTestPass)
	inst.keys = append(inst.keys, errors.JVSTestFail)
	inst.keys = append(inst.keys, errors.JVSTestWarning)
	inst.keys = append(inst.keys, errors.JVSTestUnknown)
	inst.finish = false
	return inst
}

func (s *statusCnt) done() int {
	d := 0
	for _, v := range s.cnts {
		d += v
	}
	return d
}

func (s *statusCnt) update(result *errors.JVSTestResult) {
	s.cnts[result.Status]++
}

func (s *statusCnt) StatusString() string {
	res := ""
	if s.finish {
		res += utils.Brown("[" + s.name + ":(")
	} else {
		res += utils.Brown("[" + string(s.name[0]) + ":(")
	}

	for _, k := range s.keys {
		if s.finish {
			res += errors.StatusColor(k)(errors.StatusString(k) + ":" + strconv.Itoa(s.cnts[k]))
		} else {
			res += errors.StatusColor(k)(errors.StatusShortString(k) + ":" + strconv.Itoa(s.cnts[k]))
		}
		res += utils.Brown("/")
	}
	if s.finish {
		res += utils.Brown("DONE:" + strconv.Itoa(s.done()) + "/TOTAL:" + strconv.Itoa(s.total) + ")]")
	} else {
		res += utils.Brown("D:" + strconv.Itoa(s.done()) + "/T:" + strconv.Itoa(s.total) + ")]")
	}
	return res
}

func statusMonitor(status *string, totalBuild, totalTest int, buildDone, testDone chan *errors.JVSTestResult, done chan bool) {
	buildStatus := newStatusCnt("BUILDS", totalBuild)
	testStatus := newStatusCnt("TESTS", totalTest)

LableFor:
	for {
		select {
		case result, ok := <-buildDone:
			{
				if ok {
					buildStatus.update(result)
					*status = buildStatus.StatusString() + testStatus.StatusString()
				}
				break
			}
		case result, ok := <-testDone:
			{
				if ok {
					testStatus.update(result)
					*status = buildStatus.StatusString() + testStatus.StatusString()
				}
				break
			}
		case <-done:
			buildStatus.finish = true
			testStatus.finish = true
			break LableFor
		}
	}
	*status = buildStatus.StatusString() + testStatus.StatusString()
}
