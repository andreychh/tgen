// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

// TypeNode is the concrete implementation of [TypeExpression]. Use
// [NewNamedType], [NewArrayType], or [NewUnionType] to construct a node.
type TypeNode struct {
	name     string
	inner    TypeExpression
	variants []TypeExpression
}

// NewNamedType constructs a named type node (e.g., "Integer", "Message").
func NewNamedType(name string) TypeNode {
	return TypeNode{
		name:     name,
		inner:    nil,
		variants: nil,
	}
}

// NewArrayType constructs an array type node wrapping the given element type.
func NewArrayType(inner TypeExpression) TypeNode {
	return TypeNode{
		name:     "",
		inner:    inner,
		variants: nil,
	}
}

// NewUnionType constructs a union type node from the given variant types.
func NewUnionType(variants []TypeExpression) TypeNode {
	return TypeNode{
		name:     "",
		inner:    nil,
		variants: variants,
	}
}

// Named implements [TypeExpression].
func (n TypeNode) Named() (string, bool) {
	if n.name == "" {
		return "", false
	}
	return n.name, true
}

// Array implements [TypeExpression].
//
//nolint:ireturn // TypeExpression is the intentional public contract, matching the interface definition
func (n TypeNode) Array() (TypeExpression, bool) {
	if n.inner == nil {
		return nil, false
	}
	return n.inner, true
}

// Union implements [TypeExpression].
func (n TypeNode) Union() ([]TypeExpression, bool) {
	if n.variants == nil {
		return nil, false
	}
	return n.variants, true
}

// Equal implements [TypeExpression].
func (n TypeNode) Equals(other TypeExpression) bool {
	if n.name != "" {
		name, ok := other.Named()
		return ok && n.name == name
	}
	if n.inner != nil {
		inner, ok := other.Array()
		return ok && n.inner.Equals(inner)
	}
	variants, ok := other.Union()
	if !ok || len(variants) != len(n.variants) {
		return false
	}
	for i, v := range n.variants {
		if !v.Equals(variants[i]) {
			return false
		}
	}
	return true
}
