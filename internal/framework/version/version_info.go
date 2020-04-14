package version

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"
)

type VersionInfo struct {
	GoVersion string `json:"go_version"`
	GitTag    string `json:"git_tag"`
	GitCommit string `json:"git_commit"`
	GitBranch string `json:"git_branch"`
	BuildTime string `json:"build_time"`
}

func ApikitVersion() *VersionInfo {
	return &VersionInfo{
		GoVersion: runtime.Version(),
		GitTag:    GitTag,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
		BuildTime: BuildTime,
	}
}

func (vi *VersionInfo) PrintTable() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)

	_, err := fmt.Fprintln(w, fmt.Sprintf("Go version: %s", vi.GoVersion))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Git tag: %s", vi.GitTag))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Git commit: %s", vi.GitCommit))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Git branch: %s", vi.GitBranch))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Buildtime: %s", vi.BuildTime))
	if err != nil {
		return err
	}

	return w.Flush()
}
