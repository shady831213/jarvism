package ast

import (
	"errors"
	"github.com/shady831213/jarvism/core"
	jvserrors "github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"runtime"
	"strconv"
	"strings"
)

type JVSPluginType string

const (
	JVSRunnerPlugin         = "runner"
	JVSTestDiscovererPlugin = "testDiscoverer"
	JVSSimulatorPlugin      = "simulator"
	JVSCheckerPlugin        = "checker"
)

type pluginOpts interface {
	astParser
	Name() string
}

func getPlugin(pluginType JVSPluginType, key string) pluginOpts {
	switch pluginType {
	case JVSRunnerPlugin:
		return GetRunner(key)
	case JVSSimulatorPlugin:
		return GetSimulator(key)
	case JVSTestDiscovererPlugin:
		return GetTestDiscoverer(key)
	case JVSCheckerPlugin:
		return GetChecker(key)
	}
	return nil
}

var pluginFileCache map[JVSPluginType]map[string]interface{}
var gotool string

func convertVersion(goVersion string) []int {
	s := strings.Split(strings.Replace(goVersion, "go", "", -1), ".")
	version := make([]int, 3)
	for i, v := range s {
		_v, _ := strconv.Atoi(v)
		version[i] = _v
	}
	return version
}

func compareVersion(version1, version2 string) bool {
	_version1 := convertVersion(version1)
	_version2 := convertVersion(version2)
	for i := range _version1 {
		if _version1[i] > _version2[i] {
			return true
		}
		if _version1[i] < _version2[i] {
			return false
		}
	}
	return true
}

func checkGo() error {
	if gotool != "" {
		_gotool := filepath.Join(runtime.GOROOT(), "bin", "go")
		if _, err := os.Stat(gotool); err != nil {
			if _gotool, err = exec.LookPath("go"); err != nil {
				return errors.New("can't find go tool")
			}
		}
		if !compareVersion(runtime.Version(), "go1.11.4") {
			return errors.New("go version must >= 1.11.4")
		}
		gotool = _gotool
	}
	return nil
}

func getRealPath(path string) string {
	p, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return p
}

func loadPlugin(pluginType JVSPluginType, pluginName string) *jvserrors.JVSAstError {
	//check plugin path, check customized first, then buildin
	pluginPath := path.Join(core.GetPluginsHome(), string(pluginType)+"s", pluginName)
	pluginState, err := os.Stat(pluginPath)
	if err == nil {
		symPluginPath := path.Join(core.BuildInPluginsHome(), string(pluginType)+"s", pluginName)
		defer os.RemoveAll(symPluginPath)
		if err := os.Symlink(pluginPath, symPluginPath); err != nil {
			return jvserrors.JVSPluginLoadError(pluginName, err.Error(), pluginPath)
		}
		pluginPath = symPluginPath
	} else {
		_pluginPath := path.Join(core.BuildInPluginsHome(), string(pluginType)+"s", pluginName)
		_pluginState, _err := os.Stat(_pluginPath)
		if _err != nil {
			return jvserrors.JVSPluginLoadError(pluginName, "["+err.Error()+";"+_err.Error()+"]", "["+pluginPath+";"+_pluginPath+"]")
		}
		pluginPath = _pluginPath
		pluginState = _pluginState
	}

	//check lib
	libPath := path.Join(workDir, ".jarvism_plugins", string(pluginType)+"s", pluginName+".so")
	libState, err := os.Stat(libPath)

	//compile
	if err != nil || libState.ModTime().Before(pluginState.ModTime()) {
		if err := compile(pluginType, pluginPath, libPath); err != nil {
			return jvserrors.JVSPluginLoadError(pluginName, err.Error(), getRealPath(pluginPath))
		}
	}

	if _, err := plugin.Open(libPath); err != nil {
		os.RemoveAll(path.Join(workDir, ".jarvism_plugins"))
		return jvserrors.JVSPluginLoadError(pluginName, err.Error()+" Please restart Jarvism and try again!", getRealPath(pluginPath))
	}

	return nil
}

type compileOutput struct {
	msg string
}

func (o *compileOutput) Write(p []byte) (n int, err error) {
	o.msg += string(p)
	return len(p), nil
}

func compile(pluginType JVSPluginType, pluginFile, libFile string) error {
	if err := os.MkdirAll(path.Join(workDir, string(pluginType)+"s"), os.ModePerm); err != nil {
		return err
	}
	if err := checkGo(); err != nil {
		return err
	}
	if err := os.RemoveAll(libFile); err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "-o", libFile, "-buildmode", "plugin", pluginFile)
	stderr := compileOutput{""}
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(stderr.msg + "\n" + err.Error())
	}
	return nil
}

func findAllPlugins(pluginType JVSPluginType) {
	if pluginFileCache == nil {
		pluginFileCache = make(map[JVSPluginType]map[string]interface{})
	}
	if pluginFileCache[pluginType] == nil {
		pluginFileCache[pluginType] = make(map[string]interface{})
		// custom plugin is higher priority, load buildin first then custom
		pluginFilter := func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return nil
			}
			if f.IsDir() {
				return nil
			}

			if paths := strings.Split(p, string(filepath.Separator)); paths[len(paths)-1] == "main.go" && paths[len(paths)-3] == string(pluginType)+"s" {
				if base := strings.Join(paths[:len(paths)-3], string(filepath.Separator)); base == core.BuildInPluginsHome() || base == core.GetPluginsHome() {
					pluginFileCache[pluginType][paths[len(paths)-2]] = p
				}
			}
			return nil
		}
		if err := filepath.Walk(path.Join(core.BuildInPluginsHome(), string(pluginType)+"s"), pluginFilter); err != nil {
			panic("Error in polling all plugins :" + err.Error())
		}
		if _, err := os.Stat(core.GetPluginsHome()); err == nil {
			if err := filepath.Walk(path.Join(core.GetPluginsHome(), string(pluginType)+"s"), pluginFilter); err != nil {
				panic("Error in polling all plugins :" + err.Error())
			}
		}
	}
}

func validPlugins(pluginType JVSPluginType) []string {
	findAllPlugins(pluginType)
	return utils.KeyOfStringMap(pluginFileCache[pluginType])
}
