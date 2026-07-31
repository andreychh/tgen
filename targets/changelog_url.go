// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package targets

import "github.com/andreychh/tgen/model"

// ChangelogURL represents a full URL to a release entry on the Telegram Bot API
// changelog page. A section of the API page itself is addressed by
// [TelegramURL].
type ChangelogURL struct {
	inner model.Reference
}

// NewChangelogURL creates a ChangelogURL from the reference of a release entry.
func NewChangelogURL(r model.Reference) ChangelogURL {
	return ChangelogURL{inner: r}
}

// Value returns the full URL (e.g.,
// "https://core.telegram.org/bots/api-changelog#july-14-2026").
func (u ChangelogURL) Value() string {
	return "https://core.telegram.org/bots/api-changelog#" + string(u.inner)
}
