// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package typeform models a field type as an atom repeated over a number of
// dimensions. It is the narrower of tgen's two type languages: the typeexpr
// tree can write every type the documentation composes, including a union,
// while a typeform type cannot spell a union at all.
package typeform

// Dimensionality is the number of array wrappers enclosing an atom: zero for a
// single value, one for an array of them, two for an array of arrays.
type Dimensionality int

// Type represents a field type as the atom it ultimately holds and the number
// of dimensions over which that atom repeats. The zero value holds no atom and
// is not a valid type.
type Type struct {
	atom           Atom
	dimensionality Dimensionality
}

// NewType constructs a type from its atom and dimensionality.
func NewType(atom Atom, dimensionality Dimensionality) Type {
	return Type{atom: atom, dimensionality: dimensionality}
}

// Atom returns the atom the type ultimately holds.
func (t Type) Atom() Atom {
	return t.atom
}

// Dimensionality returns the number of dimensions over which the type's atom
// repeats.
func (t Type) Dimensionality() Dimensionality {
	return t.dimensionality
}
