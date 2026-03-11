// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

import "runtime"

// LDFlagsSource provides binary metadata from values injected via -ldflags.
//
// It covers all well-known keys. Go version and platform are read from the
// runtime package since they are not available as ldflags.
type LDFlagsSource struct{}

// NewLDFlagsSource creates a LDFlagsSource.
func NewLDFlagsSource() LDFlagsSource {
	return LDFlagsSource{}
}

// Get returns the metadata value for the given key.
func (LDFlagsSource) Get(key string) (string, bool) {
	switch key {
	case KeyVersion:
		return version, true
	case KeyBuilder:
		return builder, true
	case KeyVCSRevision:
		return commit, true
	case KeyVCSTime:
		return date, true
	case KeyVCSTreeState:
		return treeState, true
	case KeyGoVersion:
		return runtime.Version(), true
	case KeyPlatform:
		return runtime.GOOS + "/" + runtime.GOARCH, true
	}
	return "", false
}
