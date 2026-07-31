// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package typebound models a field type as a resolved atom repeated over a
// number of dimensions. It is the last of tgen's three type languages: a
// typeform atom addresses a definition by reference, a typebound atom carries
// what that reference named, so a consumer renders a type without consulting a
// table. Nothing rewrites a definition after a type is bound to one, which is
// what makes carrying it safe.
package typebound

import (
	"github.com/andreychh/tgen/model/typeform"
)

// Type represents a field or return type as the resolved atom it ultimately
// holds and the number of dimensions over which that atom repeats. The zero
// value holds no atom and is not a valid type.
type Type struct {
	atom Atom
	dim  typeform.Dimensionality
}

// NewType constructs a type from its resolved atom and dimensionality.
func NewType(atom Atom, dim typeform.Dimensionality) Type {
	return Type{atom: atom, dim: dim}
}

// Atom returns the resolved atom the type ultimately holds.
func (t Type) Atom() Atom {
	return t.atom
}

// Dimensionality returns the number of dimensions over which the type's atom
// repeats.
func (t Type) Dimensionality() typeform.Dimensionality {
	return t.dim
}
