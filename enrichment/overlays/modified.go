// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/enrichment/literals"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/types"
)

// Modified represents a field whose type and optionally description are overridden.
type Modified struct {
	origin parsing.Field
	expr   types.TypeExpression
	desc   parsing.Description
}

// NewModified constructs a Modified from origin with an overridden type.
func NewModified(f parsing.Field, e types.TypeExpression) Modified {
	return NewDescribed(f, e, f.Description())
}

// NewDescribed constructs a Modified from origin with an overridden type and
// description.
func NewDescribed(f parsing.Field, e types.TypeExpression, d parsing.Description) Modified {
	return Modified{origin: f, expr: e, desc: d}
}

func (f Modified) Key() parsing.Key {
	return f.origin.Key()
}

func (f Modified) Type() parsing.Type {
	return literals.NewType(f.expr)
}

func (f Modified) Optionality() parsing.Optionality {
	return f.origin.Optionality()
}

func (f Modified) Description() parsing.Description {
	return f.desc
}
