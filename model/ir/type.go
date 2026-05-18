// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"fmt"

	"github.com/andreychh/tgen/model/types"
)

// Type represents a resolved field type expressed as a name, an array
// dimensionality, and a kind.
type Type struct {
	expr types.Expression
}

// NewType constructs a Type from a name, an array dimensionality, and
// a kind.
func NewType(expr types.Expression) Type {
	return Type{expr: expr}
}

func (t Type) Name() (string, error) { // возможно возвращать model.Name вместо string
	for expr := t.expr; ; {
		switch e := expr.(type) {
		case types.Named:
			return e.Name(), nil
		case types.Array:
			expr = e.Element()
		default:
			return "", fmt.Errorf("resolving field type: unexpected expression %q", expr)
		}
	}
}

func (t Type) Dimensionality() (int, error) {
	for expr, dim := t.expr, 0; ; {
		switch e := expr.(type) {
		case types.Named:
			return dim, nil
		case types.Array:
			expr, dim = e.Element(), dim+1
		default:
			return 0, fmt.Errorf("resolving field type: unexpected expression %q", expr)
		}
	}
}

func (t Type) Kind() (types.Kind, error) {
	for expr := t.expr; ; {
		switch e := expr.(type) {
		case types.Named:
			return e.Kind(), nil
		case types.Array:
			expr = e.Element()
		default:
			return types.KindUnknown, fmt.Errorf(
				"resolving field type: unexpected expression %q",
				expr,
			)
		}
	}
}
