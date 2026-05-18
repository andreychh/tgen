// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

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
func (r Release) Version() (string, error) {
	ver, err := r.inner.Version()
	if err != nil {
		return "", err
	}
	return string(ver), nil
}

// URL returns the URL to the release section on the Telegram Bot API page.
func (r Release) URL() (string, error) {
	ref, err := r.inner.Reference()
	if err != nil {
		return "", err
	}
	return targets.NewTelegramURL(ref).Value(), nil
}
