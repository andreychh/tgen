// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package types decodes the prose of a method's return clause into a type
// expression.
//
// [ReturnType] is the entry point: wrap a method's description prose and call
// its Value method to obtain the return type.
//
// Decoding is driven by [Rule] implementations — [ReturnsRule], [ReturnedRule],
// [ArrayRule], and [UnionRule] — each recognizing one structural form of the
// clause. The machinery they are tried by comes from [grammar]: a choice takes
// the first of them to match, and a search applies that choice at every
// position in the paragraph until one matches.
package types

import (
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/prose/grammar"
	"github.com/andreychh/tgen/model/typeexpr"
)

// Named represents an in-page anchor link, the form a return clause uses to
// identify a documented type.
type Named struct{}

// NewNamed constructs a Named.
func NewNamed() Named {
	return Named{}
}

// Matches reports whether the inline is a link to an in-page anchor and returns
// the reference it addresses. The reference is empty when the report is false.
func (Named) Matches(inline prose.Inline) (model.Reference, bool) {
	target, ok := grammar.NewAnchor().Matches(inline)
	if !ok {
		return "", false
	}
	return model.Reference(target), true
}

// Primitive represents the italic keyword form that a return clause uses for
// built-in types such as True, String, and Float.
type Primitive struct {
	primitives typeexpr.Primitives
}

// NewPrimitive constructs a Primitive over the default primitive vocabulary.
func NewPrimitive() Primitive {
	return Primitive{primitives: typeexpr.NewPrimitives()}
}

// Matches reports whether the inline is an italic text run naming a built-in
// type and returns the kind it names. Space around the word is ignored. The
// kind is empty when the report is false.
func (p Primitive) Matches(inline prose.Inline) (primitive.Kind, bool) {
	content, ok := grammar.NewItalic().Matches(inline)
	if !ok {
		return "", false
	}
	return p.primitives.Kind(strings.TrimSpace(content))
}
