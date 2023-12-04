package constant

import (
	"runtime/debug"
	"sync"
)

var (
	Version                   = ""
	coreVersion               string
	initializeCoreVersionOnce sync.Once
)

func CoreVersion() string {
	initializeCoreVersionOnce.Do(initializeCoreVersion)
	return coreVersion
}

func initializeCoreVersion() {
	if !initializeCoreVersion0() {
		coreVersion = "unknown"
	}
}

func initializeCoreVersion0() bool {
	buildInfo, loaded := debug.ReadBuildInfo()
	if !loaded {
		return false
	}
	for _, it := range buildInfo.Deps {
		if it.Path == "github.com/sagernet/sing-box" {
			coreVersion = it.Version
			return true
		}
	}
	return false
}
