// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/parsing/explicit/literals"
	"github.com/andreychh/tgen/parsing/types"
)

// Modified represents a field whose type and optionally description are overridden.
type Modified struct {
	origin explicit.Field
	expr   types.TypeExpression
	desc   explicit.Description
}

// NewModified constructs a Modified from origin with an overridden type.
func NewModified(f explicit.Field, e types.TypeExpression) Modified {
	return NewDescribed(f, e, f.Description())
}

// NewDescribed constructs a Modified from origin with an overridden type and
// description.
func NewDescribed(f explicit.Field, e types.TypeExpression, d explicit.Description) Modified {
	return Modified{origin: f, expr: e, desc: d}
}

func (f Modified) Key() explicit.Key {
	return f.origin.Key()
}

func (f Modified) Type() explicit.Type {
	return literals.NewType(f.expr)
}

func (f Modified) Optionality() explicit.Optionality {
	return f.origin.Optionality()
}

func (f Modified) Description() explicit.Description {
	return f.desc
}
