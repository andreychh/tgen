// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"strings"

	"github.com/andreychh/tgen/parsing/dom"
)

// RawObjectField is an HTML-based implementation of the Field interface. It
// interprets a single row of a Telegram API object's field table.
type RawObjectField struct {
	selection dom.Selection
}

// NewRawObjectField creates a RawObjectField from a table row (<tr>) element.
// It expects exactly three columns: Field, Type, and Description.
func NewRawObjectField(tr dom.Selection) RawObjectField {
	return RawObjectField{selection: tr}
}

// Name returns the field identifier. In this raw implementation, it defaults to
// the JSONKey, leaving case transformation to higher-level decorators.
func (f RawObjectField) Name() (string, error) {
	return f.JSONKey()
}

// Description returns the raw documentation text explaining the field's
// purpose.
func (f RawObjectField) Description() (string, error) {
	cols := f.selection.Find("td")
	if cols.Length() != 3 {
		return "", fmt.Errorf("expected 3 columns for Object field, got %d", cols.Length())
	}
	return cols.At(2).Text(), nil
}

// Type returns the raw data type string as defined in the specification.
func (f RawObjectField) Type() (string, error) {
	cols := f.selection.Find("td")
	if cols.Length() != 3 {
		return "", fmt.Errorf("expected 3 columns for Object field, got %d", cols.Length())
	}
	val := cols.At(1).Text()
	if !typeRegex.MatchString(val) {
		return "", fmt.Errorf("type %q does not match pattern %s", val, typeRegex)
	}
	return val, nil
}

// JSONKey returns the exact key used for serialization (e.g., "message_id").
func (f RawObjectField) JSONKey() (string, error) {
	cols := f.selection.Find("td")
	if cols.Length() != 3 {
		return "", fmt.Errorf("expected 3 columns for Object field, got %d", cols.Length())
	}
	val := cols.At(0).Text()
	if !jsonKeyRegex.MatchString(val) {
		return "", fmt.Errorf("json key %q does not match pattern %s", val, jsonKeyRegex)
	}
	return val, nil
}

// IsOptional reports whether the field is marked as optional in the
// documentation. It infers this by checking if the description starts with the
// "Optional" prefix.
func (f RawObjectField) IsOptional() (bool, error) {
	desc, err := f.Description()
	if err != nil {
		return false, fmt.Errorf("failed to get description for optionality check: %w", err)
	}
	return strings.HasPrefix(desc, "Optional"), nil
}
