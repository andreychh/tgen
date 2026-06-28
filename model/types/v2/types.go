// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package types models a field type as an algebraic tree of named, primitive,
// array, and union nodes. The spec's type prose maps into this tree; targets
// traverse it to render their own type syntax.
package types

// PrimitiveKind enumerates the built-in types the spec uses without an anchor
// of their own.
type PrimitiveKind string

const (
	Integer PrimitiveKind = "Integer"
	String  PrimitiveKind = "String"
	Boolean PrimitiveKind = "Boolean"
	Float   PrimitiveKind = "Float"
	True    PrimitiveKind = "True"
)

// Primitives is the vocabulary of built-in type words, mapping each spelling —
// canonical or alias — to its kind.
type Primitives struct {
	spellings map[string]PrimitiveKind
}

// NewPrimitives constructs the vocabulary of built-in type words, recognizing
// each PrimitiveKind by its canonical spelling and known aliases.
func NewPrimitives() Primitives {
	return Primitives{
		spellings: map[string]PrimitiveKind{
			"Integer": Integer,
			"Int":     Integer,
			"String":  String,
			"Boolean": Boolean,
			"Float":   Float,
			"True":    True,
		},
	}
}

// Kind returns the kind named by word, reporting whether word is a built-in type
// word in the vocabulary.
func (p Primitives) Kind(word string) (PrimitiveKind, bool) {
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
