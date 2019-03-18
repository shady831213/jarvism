package runtime

import (
	"fmt"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/utils"
	"strconv"
	"text/tabwriter"
)

type StatusCnt struct {
	Cnts  map[errors.JVSRuntimeStatus]int
	total int
	keys  []errors.JVSRuntimeStatus
	name  string
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
	res += utils.Brown("[" + string(s.name[0]) + ":(")

	for _, k := range s.keys {
		res += errors.StatusColor(k)(errors.StatusShortString(k) + ":" + strconv.Itoa(s.Cnts[k]))
		res += utils.Brown("/")
	}
	res += utils.Brown("D:" + strconv.Itoa(s.done()) + "/T:" + strconv.Itoa(s.total) + ")]")
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

func statusReport() {
	const padding = 3
	w := tabwriter.NewWriter(&stdout{}, 0, 0, padding, ' ', tabwriter.DiscardEmptyColumns|tabwriter.TabIndent|tabwriter.StripEscape|tabwriter.Debug)
	fmt.Fprintln(w, utils.Brown("Jarvism Report:"))
	fmt.Fprintln(w, " \t"+utils.Brown("TOTAL\t")+
		errors.StatusColor(errors.JVSRuntimePass)(errors.StatusString(errors.JVSRuntimePass))+"\t"+
		errors.StatusColor(errors.JVSRuntimeFail)(errors.StatusString(errors.JVSRuntimeFail))+"\t"+
		errors.StatusColor(errors.JVSRuntimeWarning)(errors.StatusString(errors.JVSRuntimeWarning))+"\t"+
		errors.StatusColor(errors.JVSRuntimeUnknown)(errors.StatusString(errors.JVSRuntimeUnknown))+"\t")
	fmt.Fprintln(w, "BUILDS\t"+utils.Brown(strconv.Itoa(buildStatus.total)+"\t")+
		errors.StatusColor(errors.JVSRuntimePass)(strconv.Itoa(buildStatus.Cnts[errors.JVSRuntimePass]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeFail)(strconv.Itoa(buildStatus.Cnts[errors.JVSRuntimeFail]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeWarning)(strconv.Itoa(buildStatus.Cnts[errors.JVSRuntimeWarning]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeUnknown)(strconv.Itoa(buildStatus.Cnts[errors.JVSRuntimeUnknown]))+"\t")
	fmt.Fprintln(w, "TESTS\t"+utils.Brown(strconv.Itoa(testStatus.total)+"\t")+
		errors.StatusColor(errors.JVSRuntimePass)(strconv.Itoa(testStatus.Cnts[errors.JVSRuntimePass]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeFail)(strconv.Itoa(testStatus.Cnts[errors.JVSRuntimeFail]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeWarning)(strconv.Itoa(testStatus.Cnts[errors.JVSRuntimeWarning]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeUnknown)(strconv.Itoa(testStatus.Cnts[errors.JVSRuntimeUnknown]))+"\t")
	fmt.Fprintln(w, utils.Brown("Jarvism Report Done!"))
	w.Flush()
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
			break LableFor
		}
	}
	statusReport()
}
