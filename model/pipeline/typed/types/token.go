// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/primitive"
)

// Token is one lexical unit of a type expression. Its variants split into value
// tokens that denote a type ([Ref], [Primitive]) and structural tokens that mark
// composition ([ArrayOf], [Or], [Series]).
//
//sumtype:decl
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

// Primitive is a built-in type word, one of the closed set of [primitive.Kind].
type Primitive struct {
	kind primitive.Kind
}

// NewPrimitive constructs a Primitive denoting the built-in type kind.
func NewPrimitive(kind primitive.Kind) Primitive {
	return Primitive{kind: kind}
}

// Kind returns the built-in type the token denotes.
func (p Primitive) Kind() primitive.Kind {
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

// Or is the word "or", dividing the alternatives of the type as a whole.
type Or struct{}

// NewOr constructs an Or token.
func NewOr() Or {
	return Or{}
}

func (Or) isToken() {}

// Series is a comma, or the "and" that closes a comma list, dividing the
// members of an enumeration. The documentation writes an array's element union
// this way and a top-level union with [Or].
type Series struct{}

// NewSeries constructs a Series token.
func NewSeries() Series {
	return Series{}
}

func (Series) isToken() {}
