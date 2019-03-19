package loader

import "github.com/shady831213/jarvism/core/plugin"

//plugin interface and add parse interface
type LoderPlugin interface {
	plugin.Plugin
	astParser
}

func getPlugin(pluginType plugin.JVSPluginType, key string) LoderPlugin {
	if p := plugin.GetPlugin(pluginType, key); p != nil {
		if v, ok := p.(LoderPlugin); ok {
			return v
		}
	}
	return nil
}
