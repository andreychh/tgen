// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/model"
)

// ReleaseURL represents the full URL to the release section on the Telegram Bot
// API page.
type ReleaseURL struct {
	inner model.Reference
}

// NewReleaseURL creates a ReleaseURL from a parsed release ref.
func NewReleaseURL(r model.Reference) ReleaseURL {
	return ReleaseURL{inner: r}
}

// AsString returns the full URL to the release section (e.g.,
// "https://core.telegram.org/bots/api#march-1-2026").
func (u ReleaseURL) AsString() (string, error) {
	ref, err := u.inner.AsString()
	if err != nil {
		return "", err
	}
	return specificationURL + "#" + ref, nil
}
