// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/targets"
)

// Release represents the latest release section of the Telegram Bot API
// specification, adapted for the Go code generation target.
type Release struct {
	inner ir.Release
}

// NewRelease creates a Release from a parsed release.
func NewRelease(r ir.Release) Release {
	return Release{inner: r}
}

// Version returns the Bot API version.
func (r Release) Version() (ReleaseVersion, error) {
	ver, err := r.inner.Version()
	if err != nil {
		return ReleaseVersion{}, err
	}
	return NewReleaseVersion(ver), nil
}

// URL returns the URL to the release section on the Telegram Bot API page.
func (r Release) URL() (targets.TelegramURL, error) {
	ref, err := r.inner.Reference()
	if err != nil {
		return targets.TelegramURL{}, err
	}
	return targets.NewTelegramURL(ref), nil
}
