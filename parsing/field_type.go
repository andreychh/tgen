// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
)

// fieldTypeRegex is compiled once at package initialization to avoid repeated
// compilation on every [FieldType.Value] call.
var fieldTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9 ,]+$`)

// FieldType represents a raw type string as it appears in the Type column of
// the Telegram Bot API field table (e.g., "Integer", "Array of String",
// "Integer or String").
//
// Call [TypeTree.Root] to parse the value into a type expression tree.
type FieldType struct {
	raw string
}

// NewFieldType creates a FieldType from a raw string.
func NewFieldType(s string) FieldType {
	return FieldType{raw: s}
}

// Value returns the validated type string. Returns an error if the string is
// not a valid field type. Structural validity is enforced during tree
// construction by [TypeTree.Root].
func (t FieldType) Value() (string, error) {
	if !fieldTypeRegex.MatchString(t.raw) {
		return "", fmt.Errorf("field type %q contains invalid characters", t.raw)
	}
	return t.raw, nil
}
