// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"strings"

	"github.com/andreychh/tgen/parsing/dom"
)

// HTMLObjectField is an HTML-based implementation of the Field interface. It
// interprets a single row of a Telegram API object's field table.
//
// Each row contains exactly three columns: Field, Type, and Description.
type HTMLObjectField struct {
	selection dom.Selection
}

// NewHTMLObjectField creates an HTMLObjectField from a table row (<tr>) element.
func NewHTMLObjectField(tr dom.Selection) HTMLObjectField {
	return HTMLObjectField{selection: tr}
}

// Key returns the field name as it appears in the JSON payload (e.g.,
// "message_id").
func (f HTMLObjectField) Key() (FieldKey, error) {
	cols, err := f.cols()
	if err != nil {
		return FieldKey{}, fmt.Errorf("reading field key: %w", err)
	}
	return NewFieldKey(cols.At(0).Text()), nil
}

// Type returns the type tree for this field. The tree is validated lazily —
// call Root() to surface structural errors.
func (f HTMLObjectField) Type() (TypeTree, error) {
	cols, err := f.cols()
	if err != nil {
		return TypeTree{}, fmt.Errorf("reading field type: %w", err)
	}
	return NewTypeTree(NewFieldType(cols.At(1).Text())), nil
}

// IsOptional reports whether the field may be absent in the JSON payload.
func (f HTMLObjectField) IsOptional() (bool, error) {
	desc, err := f.Description()
	if err != nil {
		return false, fmt.Errorf("checking field optionality: %w", err)
	}
	return strings.HasPrefix(desc, "Optional"), nil
}

// Description returns the raw documentation text from the Description column.
func (f HTMLObjectField) Description() (string, error) {
	cols, err := f.cols()
	if err != nil {
		return "", fmt.Errorf("reading field description: %w", err)
	}
	return cols.At(2).Text(), nil
}

func (f HTMLObjectField) cols() (dom.Selection, error) {
	cols := f.selection.Find("td")
	if cols.Length() != 3 {
		return nil, fmt.Errorf("expected 3 columns for object field, got %d", cols.Length())
	}
	return cols, nil
}
