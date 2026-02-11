// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package buildinfo

import (
	"runtime"
)

// Info holds metadata about the compiled binary and runtime environment.
type Info struct {
	// Version is the semantic version of the build (e.g., "v1.2.3").
	Version string
	// Commit is the full git SHA of the commit used for the build.
	Commit string
	// Branch is the name of the git branch.
	Branch string
	// BuildTime is the RFC3339 formatted timestamp of the build.
	BuildTime string
	// TreeState indicates if the git working tree was clean or dirty.
	TreeState string
	// BuiltBy identifies the tool or environment that triggered the build.
	BuiltBy string
	// GoVersion is the version of the Go compiler used.
	GoVersion string
	// Platform is the target OS and Architecture (e.g., "linux/amd64").
	Platform string
}

// ReadInfo gathers the build-time metadata and current runtime information into
// a single Info structure.
func ReadInfo() *Info {
	return &Info{
		Version:   version,
		Commit:    commit,
		Branch:    branch,
		BuildTime: buildTime,
		BuiltBy:   builtBy,
		TreeState: treeState,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}
}
