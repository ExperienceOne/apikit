package version

// Default build-time variable for lib-import.
// This file is overridden on build with build-time informations.
var (
	GitCommit string = ""
	GitBranch string = ""
	GitTag    string = ""
	BuildTime string = ""
)
