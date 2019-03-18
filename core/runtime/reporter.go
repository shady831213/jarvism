package runtime

import (
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/plugin"
)

type Reporter interface {
	plugin.Plugin
	CollectBuildResult(*errors.JVSRuntimeResult)
	CollectTestResult(*errors.JVSRuntimeResult)
	Init(totalBuild, totalTest int)
	Report()
}
