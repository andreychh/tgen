// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
)

// definitionRefRegex is compiled once at package initialization to avoid
// repeated compilation on every [DefinitionRef.Value] call.
var definitionRefRegex = regexp.MustCompile(`^[a-z0-9]+$`)

// DefinitionRef represents a reference to an Object or Method definition in the
// Telegram Bot API documentation (e.g., "message", "getupdates").
type DefinitionRef struct {
	raw string
}

// NewDefinitionRef creates a DefinitionRef from a raw string.
func NewDefinitionRef(s string) DefinitionRef {
	return DefinitionRef{raw: s}
}

// Value returns the validated reference string. Returns an error if the string
// is not a valid definition reference.
func (r DefinitionRef) Value() (string, error) {
	if !definitionRefRegex.MatchString(r.raw) {
		return "", fmt.Errorf("definition ref %q contains invalid characters", r.raw)
	}
	return r.raw, nil
}
