// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package overlays applies editorial type corrections to parsed Telegram Bot
// API fields and methods.
package overlays

import "github.com/andreychh/tgen/parsing"

// Overlay represents a conditional field transformation applied during
// enrichment.
type Overlay interface {
	Apply(f parsing.Field) parsing.Field
}
