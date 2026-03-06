// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
)

// releaseRefRegex is compiled once at package initialization to avoid repeated
// compilation on every [ReleaseRef.Value] call.
var releaseRefRegex = regexp.MustCompile(`^[a-z]+-\d+-\d+$`)

// ReleaseRef represents a reference to a release section in the Telegram Bot
// API documentation (e.g., "february-9-2026").
type ReleaseRef struct {
	raw string
}

// NewReleaseRef creates a ReleaseRef from a raw string.
func NewReleaseRef(s string) ReleaseRef {
	return ReleaseRef{raw: s}
}

// Value returns the validated reference string. Returns an error if the string
// is not a valid release reference.
func (r ReleaseRef) Value() (string, error) {
	if !releaseRefRegex.MatchString(r.raw) {
		return "", fmt.Errorf("release ref %q contains invalid characters", r.raw)
	}
	return r.raw, nil
}
