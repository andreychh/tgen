// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

import (
	"runtime"
	"runtime/debug"
)

// RuntimeSource provides binary metadata from the build information embedded by
// the Go toolchain.
//
// VCS fields (vcs.revision, vcs.time, vcs.modified) are read from
// [debug.BuildInfo.Settings]. Version and GoVersion are read from the module
// and toolchain fields respectively.
type RuntimeSource struct{}

// NewRuntimeSource creates a RuntimeSource.
func NewRuntimeSource() RuntimeSource {
	return RuntimeSource{}
}

// Get returns the metadata value for the given key. It returns "", false if
// build information is unavailable.
func (s RuntimeSource) Get(key string) (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	switch key {
	case KeyVersion:
		return info.Main.Version, true
	case KeyBuilder:
		return "gotoolchain", true
	case KeyGoVersion:
		return info.GoVersion, true
	case KeyPlatform:
		return runtime.GOOS + "/" + runtime.GOARCH, true
	}
	for _, setting := range info.Settings {
		if setting.Key == key {
			return setting.Value, true
		}
	}
	return "", false
}
