// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

import (
	"runtime"
	"runtime/debug"
)

// RuntimeSource provides binary metadata from the build information embedded
// by the Go toolchain.
//
// VCS fields (vcs.revision, vcs.time, vcs.modified) are read from
// [debug.BuildInfo.Settings]. Version and GoVersion are read from the
// module and toolchain fields respectively.
type RuntimeSource struct {
	info debug.BuildInfo
}

// NewRuntimeSource creates a RuntimeSource from the provided build information.
func NewRuntimeSource(info debug.BuildInfo) RuntimeSource {
	return RuntimeSource{info: info}
}

// Get returns the metadata value for the given key.
func (s RuntimeSource) Get(key string) (string, bool) {
	switch key {
	case KeyVersion:
		return s.info.Main.Version, true
	case KeyBuilder:
		return "gotoolchain", true
	case KeyGoVersion:
		return s.info.GoVersion, true
	case KeyPlatform:
		return runtime.GOOS + "/" + runtime.GOARCH, true
	}
	for _, setting := range s.info.Settings {
		if setting.Key == key {
			return setting.Value, true
		}
	}
	return "", false
}
