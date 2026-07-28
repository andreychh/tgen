// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package typeform

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/primitive"
)

// Atom represents a type that names what it is instead of composing it from
// other types. The concrete variants are [Named] and [Primitive].
//
//sumtype:decl
type Atom interface {
	isAtom()
}

// Named represents an atom addressing a definition by the anchor of its
// documentation.
type Named struct {
	ref model.Reference
}

// NewNamed constructs a named atom from the reference of its definition.
func NewNamed(ref model.Reference) Named {
	return Named{ref: ref}
}

// Ref returns the reference of the definition the atom addresses.
func (n Named) Ref() model.Reference {
	return n.ref
}

func (Named) isAtom() {}

// Primitive represents an atom naming a built-in type.
type Primitive struct {
	kind primitive.Kind
}

// NewPrimitive constructs a primitive atom of the given kind.
func NewPrimitive(kind primitive.Kind) Primitive {
	return Primitive{kind: kind}
}

// Kind returns the built-in type the atom names.
func (p Primitive) Kind() primitive.Kind {
	return p.kind
}

func (Primitive) isAtom() {}
