package runtime

import (
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/utils"
	"strconv"
)

type StatusCnt struct {
	Cnts   map[errors.JVSRuntimeStatus]int
	total  int
	keys   []errors.JVSRuntimeStatus
	name   string
	finish bool
}

func newStatusCnt(name string, total int) *StatusCnt {
	inst := new(StatusCnt)
	inst.name = name
	inst.Cnts = make(map[errors.JVSRuntimeStatus]int)
	inst.keys = make([]errors.JVSRuntimeStatus, 0)
	inst.total = total
	inst.Cnts[errors.JVSRuntimePass] = 0
	inst.Cnts[errors.JVSRuntimeFail] = 0
	inst.Cnts[errors.JVSRuntimeWarning] = 0
	inst.Cnts[errors.JVSRuntimeUnknown] = 0
	inst.keys = append(inst.keys, errors.JVSRuntimePass)
	inst.keys = append(inst.keys, errors.JVSRuntimeFail)
	inst.keys = append(inst.keys, errors.JVSRuntimeWarning)
	inst.keys = append(inst.keys, errors.JVSRuntimeUnknown)
	inst.finish = false
	return inst
}

func (s *StatusCnt) done() int {
	d := 0
	for _, v := range s.Cnts {
		d += v
	}
	return d
}

func (s *StatusCnt) update(result *errors.JVSRuntimeResult) {
	s.Cnts[result.Status]++
}

func (s *StatusCnt) StatusString() string {
	res := ""
	if s.finish {
		res += utils.Brown("[" + s.name + ":(")
	} else {
		res += utils.Brown("[" + string(s.name[0]) + ":(")
	}

	for _, k := range s.keys {
		if s.finish {
			res += errors.StatusColor(k)(errors.StatusString(k) + ":" + strconv.Itoa(s.Cnts[k]))
		} else {
			res += errors.StatusColor(k)(errors.StatusShortString(k) + ":" + strconv.Itoa(s.Cnts[k]))
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

var buildStatus *StatusCnt
var testStatus *StatusCnt

func GetBuildStatus() *StatusCnt {
	return buildStatus
}

func GetTestStatus() *StatusCnt {
	return testStatus
}

func statusMonitor(status *string, totalBuild, totalTest int, buildDone, testDone chan *errors.JVSRuntimeResult, done chan bool) {
	buildStatus = newStatusCnt("BUILDS", totalBuild)
	testStatus = newStatusCnt("TESTS", totalTest)
	*status = buildStatus.StatusString() + testStatus.StatusString()

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
