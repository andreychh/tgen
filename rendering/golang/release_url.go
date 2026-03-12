// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/parsing"

// ReleaseURL represents the full URL to the release section on the Telegram Bot
// API page.
type ReleaseURL struct {
	inner parsing.ReleaseRef
}

// NewReleaseURL creates a ReleaseURL from a parsed release ref.
func NewReleaseURL(r parsing.ReleaseRef) ReleaseURL {
	return ReleaseURL{inner: r}
}

// Value returns the full URL to the release section (e.g.,
// "https://core.telegram.org/bots/api#march-1-2026").
func (u ReleaseURL) Value() (string, error) {
	ref, err := u.inner.Value()
	if err != nil {
		return "", err
	}
	return specificationURL + "#" + ref, nil
}
