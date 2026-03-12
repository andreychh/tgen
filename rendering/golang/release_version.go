// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/parsing"

// ReleaseVersion represents the Bot API version of the specification, adapted
// for the Go code generation target.
type ReleaseVersion struct {
	inner parsing.ReleaseVersion
}

// NewReleaseVersion creates a ReleaseVersion from a parsed release version.
func NewReleaseVersion(v parsing.ReleaseVersion) ReleaseVersion {
	return ReleaseVersion{inner: v}
}

// Value returns the Bot API version string (e.g., "9.5").
func (v ReleaseVersion) Value() (string, error) {
	return v.inner.Value()
}
