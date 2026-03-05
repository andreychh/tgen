// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

// TypeExpression represents a node in the type expression tree parsed from the
// Telegram Bot API field table. A node is one of three kinds: a named type, an
// array type, or a union type.
//
// Exactly one of the three methods returns a non-zero value for any given node.
// The rendering layer traverses the tree by calling each method and acting on
// the one that returns true.
type TypeExpression interface {
	// Named returns the type name and true if this node is a named type (e.g.,
	// "Integer", "Message").
	Named() (string, bool)

	// Array returns the element type and true if this node is an array type (e.g.,
	// "Array of Integer").
	Array() (TypeExpression, bool)

	// Union returns the variant types and true if this node is a union type (e.g.,
	// "Integer or String").
	Union() ([]TypeExpression, bool)
}

// TypeNode is the concrete implementation of [TypeExpression]. Use
// [NewNamedType], [NewArrayType], or [NewUnionType] to construct a node.
type TypeNode struct {
	name     string
	inner    TypeExpression
	variants []TypeExpression
}

// NewNamedType constructs a named type node (e.g., "Integer", "Message").
func NewNamedType(name string) TypeNode {
	return TypeNode{name: name}
}

// NewArrayType constructs an array type node wrapping the given element type.
func NewArrayType(inner TypeExpression) TypeNode {
	return TypeNode{inner: inner}
}

// NewUnionType constructs a union type node from the given variant types.
func NewUnionType(variants []TypeExpression) TypeNode {
	return TypeNode{variants: variants}
}

func (n TypeNode) Named() (string, bool) {
	if n.name == "" {
		return "", false
	}
	return n.name, true
}

func (n TypeNode) Array() (TypeExpression, bool) {
	if n.inner == nil {
		return nil, false
	}
	return n.inner, true
}

func (n TypeNode) Union() ([]TypeExpression, bool) {
	if n.variants == nil {
		return nil, false
	}
	return n.variants, true
}
