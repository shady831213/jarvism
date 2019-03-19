package core

/*
define some global path
*/
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

func getCfgFile(p string) (string, error) {
	if s, err := os.Stat(path.Join(p, "jarvism_cfg")); err == nil && s.IsDir() {
		return path.Join(p, "jarvism_cfg"), os.Setenv("JVS_PRJ_HOME", p)
	}
	if s, err := os.Stat(path.Join(p, "jarvism_cfg.yaml")); err == nil && !s.IsDir() {
		return path.Join(p, "jarvism_cfg.yaml"), os.Setenv("JVS_PRJ_HOME", p)
	}
	return "", errors.New("no \"jarvism_cfg\" dir or \"jarvism_cfg.yaml\" found in " + p + "!")
}

func GetCfgFile() (string, error) {
	if os.Getenv("JVS_PRJ_HOME") == "" {
		pwd, err := filepath.Abs(os.Getenv("PWD"))
		if err != nil {
			return "", errors.New("Env $JVS_PRJ_HOME is not set and no \"jarvism_cfg\" dir or \"jarvism_cfg.yaml\" found in current path tree!" + err.Error())
		}
		for p := pwd; p != string(filepath.Separator); p = filepath.Dir(p) {
			if cfgFile, err := getCfgFile(p); cfgFile != "" {
				return cfgFile, err
			}
		}
		return "", errors.New("Env $JVS_PRJ_HOME is not set and no \"jarvism_cfg\" dir or \"jarvism_cfg.yaml\" found in current path tree!")
	}
	p, err := filepath.Abs(os.ExpandEnv(os.Getenv("JVS_PRJ_HOME")))
	if err != nil {
		return "", errors.New("no \"jarvism_cfg\" dir or \"jarvism_cfg.yaml\" found in  $JVS_PRJ_HOME(" + os.ExpandEnv("JVS_PRJ_HOME") + ")!" + err.Error())
	}
	return getCfgFile(p)
}

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
