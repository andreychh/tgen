// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/targets"
)

// Release represents the Bot API release the generated files were read from.
type Release struct {
	inner ir.Release
}

// NewRelease creates a Release from the record of a release.
func NewRelease(r ir.Release) Release {
	return Release{inner: r}
}

// Version returns the Bot API version of the release.
func (r Release) Version() string {
	return string(r.inner.Version)
}

// Changelog returns the URL of the changelog entry announcing the release.
func (r Release) Changelog() string {
	return targets.NewChangelogURL(r.inner.Ref).Value()
}
