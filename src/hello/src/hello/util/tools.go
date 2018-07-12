package util

import (
	"path/filepath"
	"os"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"time"
	"net"
)

const (
	ConfigFilePath = "conf"
)

func GetCfgFilePath() string {
	var err error
	appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", appPath)
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", workPath)
	appConfigPath := filepath.Join(workPath, ConfigFilePath)
	if !FileExists(appConfigPath) {
		appConfigPath = filepath.Join(appPath, ConfigFilePath)
		if !FileExists(appConfigPath) {
			goPath := GetGoPath()
			for _, val := range goPath {
				appConfigPath = filepath.Join(val, "src", "appmgr", ConfigFilePath)
				fmt.Println(appConfigPath)
				if FileExists(appConfigPath) {
					return appConfigPath
				}
			}
			appConfigPath = "/"
		}
	}
	return appConfigPath
}

func GetGoPath() []string {
	goPath := os.Getenv("GOPATH")
	fmt.Println(goPath)
	if strings.Contains(goPath, ";") { //windows
		return strings.Split(goPath, ";")
	} else if strings.Contains(goPath, ":") { //linux
		return strings.Split(goPath, ":")
	} else { //only one
		path := make([]string, 1, 1)
		path[0] = goPath
		return path
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
