// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
)

// fieldKeyRegex is compiled once at package initialization to avoid repeated
// compilation on every [FieldKey.Value] call.
var fieldKeyRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

// FieldKey represents a JSON property key as it appears in the Field column of
// the Telegram Bot API field table (e.g., "message_id", "first_name").
//
// Call [FieldKey.Value] to retrieve the validated string.
type FieldKey struct {
	raw string
}

// NewFieldKey creates a FieldKey from a raw string.
func NewFieldKey(s string) FieldKey {
	return FieldKey{raw: s}
}

// Value returns the validated JSON key string. Returns an error if the raw
// string is not a valid JSON key.
func (k FieldKey) Value() (string, error) {
	if !fieldKeyRegex.MatchString(k.raw) {
		return "", fmt.Errorf("invalid field key: %q", k.raw)
	}
	return k.raw, nil
}
