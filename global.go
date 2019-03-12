package jarivsm

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var pkgPath string

func PkgPath() string {
	return getPkgPath()
}

func CorePath() string {
	return path.Join(getPkgPath(), "core")
}

func BuildInOptionPath() string {
	return path.Join(getPkgPath(), "buildInOptions")
}

func TestDiscoverersPath() string {
	return path.Join(getPkgPath(), "testDiscoverers")
}

func simulatorsPath() string {
	return path.Join(getPkgPath(), "simulators")
}

func getPkgPath() string {
	if pkgPath == "" {
		_, file, _, ok := runtime.Caller(1)
		if !ok {
			panic(errors.New("Can not get current file info"))
		}
		_pkgPath, err := filepath.Abs(path.Dir(file))
		if err != nil {
			panic(err)
		}
		pkgPath = _pkgPath
	}
	return pkgPath
}

var prjHome string

func CheckEnv() error {
	if os.Getenv("JVS_PRJ_HOME") == "" {
		return errors.New("Env $JVS_PRJ_HOME is not set!")
	}
	prjHome = os.Getenv("JVS_PRJ_HOME")
	return nil
}

func GetPrjHome() string {
	return prjHome
}
