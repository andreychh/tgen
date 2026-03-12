// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/parsing"

// Release represents the latest release section of the Telegram Bot API
// specification, adapted for the Go code generation target.
type Release struct {
	inner parsing.Release
}

// NewRelease creates a Release from a parsed release.
func NewRelease(r parsing.Release) Release {
	return Release{inner: r}
}

// Version returns the Bot API version.
func (r Release) Version() ReleaseVersion {
	return NewReleaseVersion(r.inner.Version())
}

// URL returns the URL to the release section on the Telegram Bot API page.
func (r Release) URL() ReleaseURL {
	return NewReleaseURL(r.inner.Ref())
}
