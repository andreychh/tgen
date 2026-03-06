// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
)

// objectNameRegex is compiled once at package initialization to avoid repeated
// compilation on every [ObjectName.Value] call.
var objectNameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`)

// ObjectName represents the name of an Object or Union in the Telegram Bot API
// documentation (e.g., "Message", "MaybeInaccessibleMessage").
type ObjectName struct {
	raw string
}

// NewObjectName creates an ObjectName from a raw string.
func NewObjectName(s string) ObjectName {
	return ObjectName{raw: s}
}

// Value returns the validated object name string. Returns an error if the
// string is not a valid object name.
func (n ObjectName) Value() (string, error) {
	if !objectNameRegex.MatchString(n.raw) {
		return "", fmt.Errorf("object name %q contains invalid characters", n.raw)
	}
	return n.raw, nil
}
