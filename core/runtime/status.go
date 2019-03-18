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

type statusReporter struct {
	buildStatus *StatusCnt
	testStatus  *StatusCnt
	status      string
}

func (r *statusReporter) Name() string {
	return "statusReporter"
}

func (r *statusReporter) Init(totalBuild, totalTest int) {
	r.buildStatus = newStatusCnt("BUILDS", totalBuild)
	r.testStatus = newStatusCnt("TESTS", totalTest)
	r.status = r.buildStatus.StatusString() + r.testStatus.StatusString()
}

func (r *statusReporter) CollectBuildResult(result *errors.JVSRuntimeResult) {
	r.buildStatus.update(result)
	r.status = r.buildStatus.StatusString() + r.testStatus.StatusString()
}

func (r *statusReporter) CollectTestResult(result *errors.JVSRuntimeResult) {
	r.testStatus.update(result)
	r.status = r.buildStatus.StatusString() + r.testStatus.StatusString()
}

func (r *statusReporter) Report() {
	const padding = 3
	w := tabwriter.NewWriter(&stdout{}, 0, 0, padding, ' ', tabwriter.DiscardEmptyColumns|tabwriter.TabIndent|tabwriter.StripEscape|tabwriter.Debug)
	fmt.Fprintln(w, utils.Brown("Jarvism Report:"))
	fmt.Fprintln(w, " \t"+utils.Brown("TOTAL\t")+
		errors.StatusColor(errors.JVSRuntimePass)(errors.StatusString(errors.JVSRuntimePass))+"\t"+
		errors.StatusColor(errors.JVSRuntimeFail)(errors.StatusString(errors.JVSRuntimeFail))+"\t"+
		errors.StatusColor(errors.JVSRuntimeWarning)(errors.StatusString(errors.JVSRuntimeWarning))+"\t"+
		errors.StatusColor(errors.JVSRuntimeUnknown)(errors.StatusString(errors.JVSRuntimeUnknown))+"\t")
	fmt.Fprintln(w, "BUILDS\t"+utils.Brown(strconv.Itoa(r.buildStatus.total)+"\t")+
		errors.StatusColor(errors.JVSRuntimePass)(strconv.Itoa(r.buildStatus.Cnts[errors.JVSRuntimePass]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeFail)(strconv.Itoa(r.buildStatus.Cnts[errors.JVSRuntimeFail]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeWarning)(strconv.Itoa(r.buildStatus.Cnts[errors.JVSRuntimeWarning]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeUnknown)(strconv.Itoa(r.buildStatus.Cnts[errors.JVSRuntimeUnknown]))+"\t")
	fmt.Fprintln(w, "TESTS\t"+utils.Brown(strconv.Itoa(r.testStatus.total)+"\t")+
		errors.StatusColor(errors.JVSRuntimePass)(strconv.Itoa(r.testStatus.Cnts[errors.JVSRuntimePass]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeFail)(strconv.Itoa(r.testStatus.Cnts[errors.JVSRuntimeFail]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeWarning)(strconv.Itoa(r.testStatus.Cnts[errors.JVSRuntimeWarning]))+"\t"+
		errors.StatusColor(errors.JVSRuntimeUnknown)(strconv.Itoa(r.testStatus.Cnts[errors.JVSRuntimeUnknown]))+"\t")
	fmt.Fprintln(w, utils.Brown("Jarvism Report Done!"))
	w.Flush()
}

var status = &statusReporter{}

func GetBuildStatus() *StatusCnt {
	return status.buildStatus
}
func GetTestStatus() *StatusCnt {
	return status.testStatus
}
