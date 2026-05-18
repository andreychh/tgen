// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
)

// Release represents the release metadata of the Telegram Bot API.
type Release struct {
	inner spec.Release
}

// NewRelease constructs a Release from a parsed release.
func NewRelease(r spec.Release) Release {
	return Release{inner: r}
}

func (r Release) Reference() (model.Reference, error) {
	return r.inner.Reference()
}

func (r Release) Version() (model.ReleaseVersion, error) {
	return r.inner.Version()
}

func (r Release) Date() (model.ReleaseDate, error) {
	return r.inner.Date()
}
