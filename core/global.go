package core

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var pkgPath string

func PkgPath() string {
	return getJarvismPath()
}

func BuildInPluginsHome() string {
	return path.Join(getJarvismPath(), "plugins")
}

func CorePath() string {
	return path.Join(getJarvismPath(), "core")
}

func BuildInOptionPath() string {
	return path.Join(CorePath(), "buildInOptions")
}

func TestDiscoverersPath() string {
	return path.Join(BuildInPluginsHome(), "testDiscoverers")
}

func SimulatorsPath() string {
	return path.Join(BuildInPluginsHome(), "simulators")
}

func RunnersPath() string {
	return path.Join(BuildInPluginsHome(), "runners")
}

func getJarvismPath() string {
	if pkgPath == "" {
		_, file, _, ok := runtime.Caller(1)
		if !ok {
			panic(errors.New("Can not get current file info"))
		}
		_pkgPath, err := filepath.Abs(path.Dir(path.Dir(file)))
		if err != nil {
			panic(err)
		}
		pkgPath = _pkgPath
	}
	return pkgPath
}

var prjHome string
var pluginsHome string

func CheckEnv() error {
	if os.Getenv("JVS_PRJ_HOME") == "" {
		return errors.New("Env $JVS_PRJ_HOME is not set!")
	}
	prjHome = os.ExpandEnv(os.Getenv("JVS_PRJ_HOME"))
	if os.Getenv("JVS_PLUGINS_HOME") == "" {
		pluginsHome = path.Join(prjHome, "jarvism_plugins")
		return nil
	}
	pluginsHome = os.ExpandEnv(os.Getenv("JVS_PLUGINS_HOME"))
	return nil
}

func GetPrjHome() string {
	return prjHome
}

func GetPluginsHome() string {
	return pluginsHome
}
