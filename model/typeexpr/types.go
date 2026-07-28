// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package typeexpr models a field type as an algebraic tree of named,
// primitive, array, and union nodes. The type prose of the documentation maps
// into this tree, the wider of tgen's two type languages: it can write every
// type the documentation composes, including a union.
package typeexpr

import (
	"github.com/andreychh/tgen/model/primitive"
)

// Primitives is the vocabulary of built-in type words, mapping each spelling —
// canonical or alias — to its kind.
type Primitives struct {
	spellings map[string]primitive.Kind
}

// NewPrimitives constructs the vocabulary of built-in type words, recognizing
// each [primitive.Kind] by its canonical spelling and known aliases.
func NewPrimitives() Primitives {
	return Primitives{
		spellings: map[string]primitive.Kind{
			"Integer": primitive.Integer,
			"Int":     primitive.Integer,
			"String":  primitive.String,
			"Boolean": primitive.Boolean,
			"Float":   primitive.Float,
			"True":    primitive.True,
		},
	}
}

// Kind returns the kind named by word, reporting whether word is a built-in type
// word in the vocabulary.
func (p Primitives) Kind(word string) (primitive.Kind, bool) {
	kind, ok := p.spellings[word]
	return kind, ok
}

// Expression represents a field type as a tree. The concrete variants are
// Named, Primitive, Array, and Union.
//
//sumtype:decl
type Expression interface {
	Equals(other Expression) bool
	isNode()
}
