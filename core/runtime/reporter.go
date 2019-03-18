package runtime

import (
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/plugin"
)

type Reporter interface {
	plugin.Plugin
	CollectBuildResult(*errors.JVSRuntimeResult)
	CollectTestResult(*errors.JVSRuntimeResult)
	Init(jobId string, totalBuild, totalTest int)
	Report()
}

func RegisterReporter(c func() plugin.Plugin) {
	plugin.RegisterPlugin(plugin.JVSReporterPlugin, c)
}
