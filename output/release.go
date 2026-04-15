// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package output

import "github.com/andreychh/tgen/meta"

// Release represents the release metadata of the tool, adapted for use in
// templates.
type Release struct {
	inner meta.Release
}

// NewRelease creates a Release from the given [meta.Release].
func NewRelease(r meta.Release) Release {
	return Release{inner: r}
}

// Version returns the semantic version of the tool (e.g., "v0.1.0").
func (r Release) Version() string {
	return r.inner.Version()
}
