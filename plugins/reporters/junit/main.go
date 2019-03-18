package main

import (
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/options"
	"github.com/shady831213/jarvism/core/plugin"
	"github.com/shady831213/jarvism/core/runtime"
	"github.com/shady831213/jarvism/core/utils"
	"os"
	"path"
)

type junitReporter struct {
	jobId string
}

func newJunitReporter() plugin.Plugin {
	return new(junitReporter)
}

func (r *junitReporter) Name() string {
	return "junit"
}

func (r *junitReporter) Init(jobId string, totalBuild, totalTest int) {
	r.jobId = jobId
	initJunitXml(totalBuild, totalTest)
}

func (r *junitReporter) CollectBuildResult(result *errors.JVSRuntimeResult) {
	updateBuild(result)
}

func (r *junitReporter) CollectTestResult(result *errors.JVSRuntimeResult) {
	updateTest(result)
}

func (r *junitReporter) Report() {
	file := path.Join(core.GetReportDir(), r.Name(), r.jobId+".xml")
	if err := os.MkdirAll(path.Join(core.GetReportDir(), r.Name()), os.ModePerm); err != nil {
		runtime.Println(utils.LightRed("gen junit report " + file + " failed!\n" + err.Error()))
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	if err != nil {
		runtime.Println(utils.LightRed("gen junit report " + file + " failed!\n" + err.Error()))
	}
	if err := writeReport(f); err != nil {
		runtime.Println(utils.LightRed("gen junit report " + file + " failed!\n" + err.Error()))
	}
	runtime.Println(utils.Brown("gen junit report " + file + "!"))
}

var noXMLHeader bool

func init() {
	runtime.RegisterReporter(newJunitReporter)
	options.GetJvsOptions().BoolVar(&noXMLHeader, "junitNoXMLHeader", false, "if enable, xmlHeader will not be generated")
}
