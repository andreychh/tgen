// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/model/types"
)

// Modified represents a field whose type and optionally description are overridden.
type Modified struct {
	origin explicit.Field
	expr   types.Expression
	desc   model.Description
}

// NewModified constructs a Modified from origin with an overridden type.
func NewModified(f explicit.Field, e types.Expression) Modified {
	return NewDescribed(f, e, f.Description())
}

// NewDescribed constructs a Modified from origin with an overridden type and
// description.
func NewDescribed(f explicit.Field, e types.Expression, d model.Description) Modified {
	return Modified{origin: f, expr: e, desc: d}
}

func (f Modified) Key() (model.Key, error) {
	return f.origin.Key()
}

func (f Modified) Type() (types.Expression, error) {
	return f.expr, nil
}

func (f Modified) Optionality() (model.Optionality, error) {
	return f.origin.Optionality()
}

func (f Modified) Description() model.Description {
	return f.desc
}
