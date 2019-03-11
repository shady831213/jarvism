package core

import (
	"github.com/shady831213/jarvism/utils"
	"strconv"
)

type JVSTestStatus int

const (
	_ JVSTestStatus = iota
	JVSTestPass
	JVSTestWarning
	JVSTestFail
	JVSTestUnknown
)

func StatusColor(status JVSTestStatus) func(str string, modifier ...interface{}) string {
	switch status {
	case JVSTestPass:
		return utils.Green
	case JVSTestWarning:
		return utils.Yellow
	case JVSTestFail:
		return utils.Red
	case JVSTestUnknown:
		return utils.LightRed
	}
	return utils.LightRed
}

func StatusString(status JVSTestStatus) string {
	switch status {
	case JVSTestPass:
		return "PASS"
	case JVSTestWarning:
		return "WARNING"
	case JVSTestFail:
		return "FAIL"
	case JVSTestUnknown:
		return "UNKNOWN"
	}
	return "UNKNOWN"
}

func StatusShortString(status JVSTestStatus) string {
	switch status {
	case JVSTestPass:
		return "P"
	case JVSTestWarning:
		return "W"
	case JVSTestFail:
		return "F"
	case JVSTestUnknown:
		return "U"
	}
	return "U"
}

type JVSTestResult struct {
	status JVSTestStatus
	msg    string
}

func (e *JVSTestResult) Error() string {
	return StatusColor(e.status)(StatusString(e.status) + e.msg)
}

func JVSTestResultPass(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestPass
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + msg
	}
	return inst
}

func JVSTestResultFail(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestFail
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + "Error:" + msg
	}
	return inst
}

func JVSTestResultWarning(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestWarning
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + "Warning:" + msg
	}
	return inst
}

func JVSTestResultUnknown(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestUnknown
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + "UnKnown:" + msg
	}
	return inst
}

type statusCnt struct {
	cnts   map[JVSTestStatus]int
	total  int
	keys   []JVSTestStatus
	name   string
	finish bool
}

func newStatusCnt(name string, total int) *statusCnt {
	inst := new(statusCnt)
	inst.name = name
	inst.cnts = make(map[JVSTestStatus]int)
	inst.keys = make([]JVSTestStatus, 0)
	inst.total = total
	inst.cnts[JVSTestPass] = 0
	inst.cnts[JVSTestFail] = 0
	inst.cnts[JVSTestWarning] = 0
	inst.cnts[JVSTestUnknown] = 0
	inst.keys = append(inst.keys, JVSTestPass)
	inst.keys = append(inst.keys, JVSTestFail)
	inst.keys = append(inst.keys, JVSTestWarning)
	inst.keys = append(inst.keys, JVSTestUnknown)
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

func (s *statusCnt) update(result *JVSTestResult) {
	s.cnts[result.status]++
}

func (s *statusCnt) statusString() string {
	res := ""
	if s.finish {
		res += utils.Brown("[" + s.name + ":(")
	} else {
		res += utils.Brown("[" + string(s.name[0]) + ":(")
	}

	for _, k := range s.keys {
		if s.finish {
			res += StatusColor(k)(StatusString(k) + ":" + strconv.Itoa(s.cnts[k]))
		} else {
			res += StatusColor(k)(StatusShortString(k) + ":" + strconv.Itoa(s.cnts[k]))
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

func statusMonitor(status *string, totalBuild, totalTest int, buildDone, testDone chan *JVSTestResult, done chan bool) {
	buildStatus := newStatusCnt("BUILDS", totalBuild)
	testStatus := newStatusCnt("TESTS", totalTest)

LableFor:
	for {
		select {
		case result, ok := <-buildDone:
			{
				if ok {
					buildStatus.update(result)
					*status = buildStatus.statusString() + testStatus.statusString()
				}
				break
			}
		case result, ok := <-testDone:
			{
				if ok {
					testStatus.update(result)
					*status = buildStatus.statusString() + testStatus.statusString()
				}
				break
			}
		case <-done:
			buildStatus.finish = true
			testStatus.finish = true
			break LableFor
		}
	}
	*status = buildStatus.statusString() + testStatus.statusString()
}
