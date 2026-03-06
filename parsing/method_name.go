// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
)

// methodNameRegex is compiled once at package initialization to avoid repeated
// compilation on every [MethodName.Value] call.
var methodNameRegex = regexp.MustCompile(`^[a-z][a-zA-Z0-9]+$`)

// MethodName represents the name of a Method in the Telegram Bot API
// documentation (e.g., "sendMessage", "getUpdates").
type MethodName struct {
	raw string
}

// NewMethodName creates a MethodName from a raw string.
func NewMethodName(s string) MethodName {
	return MethodName{raw: s}
}

// Value returns the validated method name string. Returns an error if the
// string is not a valid method name.
func (n MethodName) Value() (string, error) {
	if !methodNameRegex.MatchString(n.raw) {
		return "", fmt.Errorf("method name %q contains invalid characters", n.raw)
	}
	return n.raw, nil
}
