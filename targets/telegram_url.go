// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package targets

import "github.com/andreychh/tgen/model"

// TelegramURL represents a full URL to a section on the Telegram Bot API page.
type TelegramURL struct {
	inner model.Reference
}

// NewTelegramURL creates a TelegramURL from a parsed reference.
func NewTelegramURL(r model.Reference) TelegramURL {
	return TelegramURL{inner: r}
}

// AsString returns the full URL (e.g.,
// "https://core.telegram.org/bots/api#march-1-2026").
func (u TelegramURL) AsString() (string, error) {
	ref, err := u.inner.AsString()
	if err != nil {
		return "", err
	}
	return "https://core.telegram.org/bots/api#" + ref, nil
}
