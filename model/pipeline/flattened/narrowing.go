// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package flattened

import (
	"errors"

	"github.com/andreychh/tgen/model/typeexpr"
	"github.com/andreychh/tgen/model/typeform"
)

// Narrowing is a type expression ready to be reduced to the atom it holds and
// the number of dimensions over which that atom repeats.
type Narrowing struct {
	expr typeexpr.Expression
}

// NewNarrowing constructs a Narrowing over a type expression.
func NewNarrowing(expr typeexpr.Expression) Narrowing {
	return Narrowing{expr: expr}
}

// Value returns the flat type the expression denotes. It fails when the
// expression holds a union, which has no flat form.
func (n Narrowing) Value() (typeform.Type, error) {
	expr, dim := n.expr, typeform.Dimensionality(0)
	for {
		array, ok := expr.(typeexpr.Array)
		if !ok {
			break
		}
		expr, dim = array.Element(), dim+1
	}
	atom, err := n.atom(expr)
	if err != nil {
		return typeform.Type{}, err
	}
	return typeform.NewType(atom, dim), nil
}

// atom returns the atom expr denotes. It fails on any expression that composes
// a type out of others rather than naming one.
func (n Narrowing) atom(expr typeexpr.Expression) (typeform.Atom, error) {
	switch expr := expr.(type) {
	case typeexpr.Named:
		return typeform.NewNamed(expr.Ref()), nil
	case typeexpr.Primitive:
		return typeform.NewPrimitive(expr.Kind()), nil
	case typeexpr.Array:
		return nil, errors.New("array left below the counted dimensions")
	case typeexpr.Union:
		return nil, errors.New("union has no flat form")
	}
	return nil, errors.New("unknown type expression")
}
