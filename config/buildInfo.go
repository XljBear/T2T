package config

import "runtime"

var (
	CommitID  string = "unknown"      // git commit ID
	Version          = "v0.0.1"       // version number
	BuildTime string = "unknown"      // build time
	GoVersion string = "unknown"      // go version
	OS               = runtime.GOOS   // operating system
	Arch             = runtime.GOARCH // architecture
)
