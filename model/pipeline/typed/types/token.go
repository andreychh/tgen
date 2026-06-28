// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types/v2"
)

// Token is one lexical unit of a type expression. Its variants split into value
// tokens that denote a type ([Ref], [Primitive]) and structural tokens that mark
// composition ([ArrayOf], [Separator]).
type Token interface {
	isToken()
}

// Ref is a reference to a documented type, decoded from a link's anchor.
type Ref struct {
	reference model.Reference
}

// NewRef constructs a Ref denoting the documented type at reference.
func NewRef(reference model.Reference) Ref {
	return Ref{reference: reference}
}

// Reference returns the reference to the documented type the token denotes.
func (r Ref) Reference() model.Reference {
	return r.reference
}

func (Ref) isToken() {}

// Primitive is a built-in type word, one of the closed PrimitiveKind set.
type Primitive struct {
	kind types.PrimitiveKind
}

// NewPrimitive constructs a Primitive denoting the built-in type kind.
func NewPrimitive(kind types.PrimitiveKind) Primitive {
	return Primitive{kind: kind}
}

// Kind returns the built-in type the token denotes.
func (p Primitive) Kind() types.PrimitiveKind {
	return p.kind
}

func (Primitive) isToken() {}

// ArrayOf is the "Array of" prefix that wraps the following type.
type ArrayOf struct{}

// NewArrayOf constructs an ArrayOf prefix token.
func NewArrayOf() ArrayOf {
	return ArrayOf{}
}

func (ArrayOf) isToken() {}

// Separator divides the variants of a union. The spec writes it as "or" at the
// top level and as commas with a final "and" inside an array element; all three
// mean the same thing.
type Separator struct{}

// NewSeparator constructs a Separator token.
func NewSeparator() Separator {
	return Separator{}
}

func (Separator) isToken() {}
