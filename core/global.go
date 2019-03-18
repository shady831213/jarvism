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
var workDir string

func CheckEnv() error {
	if os.Getenv("JVS_PRJ_HOME") == "" {
		return errors.New("Env $JVS_PRJ_HOME is not set!")
	}
	p, err := filepath.Abs(os.ExpandEnv(os.Getenv("JVS_PRJ_HOME")))
	if err != nil {
		return err
	}
	prjHome = p

	pluginsHome = os.ExpandEnv(os.Getenv("JVS_PLUGINS_HOME"))
	if pluginsHome == "" {
		pluginsHome = path.Join(prjHome, "jarvism_plugins")
	}

	workDir = os.ExpandEnv(os.Getenv("JVS_WORK_DIR"))
	if workDir == "" {
		workDir = path.Join(prjHome, "work")
	}

	if err := os.MkdirAll(workDir, os.ModePerm); err != nil {
		return err
	}
	if _, err := os.Stat(workDir); err != nil {
		return err
	}
	return nil
}

func GetPrjHome() string {
	return prjHome
}

func GetPluginsHome() string {
	return pluginsHome
}

func GetWorkDir() string {
	return workDir
}

func GetReportDir() string {
	return path.Join(GetWorkDir(), "report")
}
