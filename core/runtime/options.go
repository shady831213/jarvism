package runtime

import (
	"github.com/shady831213/jarvism/core/options"
	"github.com/shady831213/jarvism/core/plugin"
	"github.com/shady831213/jarvism/core/utils"
	"strings"
)

var runTimeMaxJob int
var runTimeSimOnly bool
var runTimeUnique bool
var runTimeReporter = &runTimeReporterVar{}

type runTimeReporterVar struct {
	reporters *utils.StringMapSet
}

func (v *runTimeReporterVar) getReporters() []Reporter {
	reporters := make([]Reporter, 0)
	for _, r := range v.reporters.List() {
		reporters = append(reporters, r.(Reporter))
	}
	return reporters
}

func (v *runTimeReporterVar) Set(s string) error {
	if v.reporters == nil {
		v.reporters = utils.NewStringMapSet()
	}
	reporter := plugin.GetPlugin(plugin.JVSReporterPlugin, s)
	if reporter == nil {
		if err := plugin.LoadPlugin(plugin.JVSReporterPlugin, s); err != nil {
			return err
		}
		reporter = plugin.GetPlugin(plugin.JVSReporterPlugin, s)
	}
	v.reporters.Add(s, reporter)
	return nil
}

func (v *runTimeReporterVar) String() string {
	if v.reporters == nil {
		v.reporters = utils.NewStringMapSet()
	}
	return strings.Join(v.reporters.Keys(), " ")
}

func (v *runTimeReporterVar) IsBoolFlag() bool {
	return false
}

func init() {
	options.GetJvsOptions().IntVar(&runTimeMaxJob, "max_job", -1, "limit of runtime coroutines, default is unlimited.")
	options.GetJvsOptions().BoolVar(&runTimeSimOnly, "sim_only", false, "bypass compile and only run simulation, default is false.")
	options.GetJvsOptions().BoolVar(&runTimeUnique, "unique", false, "if set jobId(timestamp) will be included in hash, then builds and testcases will have unique name and be in unique dir.default is false.")
	options.GetJvsOptions().Var(runTimeReporter, "reporter", "add reporter plugin, can apply multi times, default")
}
