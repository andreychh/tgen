// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"slices"
)

// Union represents a choice between several variant types.
type Union struct {
	variants []Expression
}

// NewUnion constructs a union type from its variant types.
func NewUnion(variants ...Expression) Union {
	return Union{variants: variants}
}

// Variants returns the variant types of the union.
func (u Union) Variants() []Expression {
	return u.variants
}

// Equals reports whether other is a Union with equal variant types in order.
func (u Union) Equals(other Expression) bool {
	if other, ok := other.(Union); ok {
		return slices.EqualFunc(u.variants, other.variants, Expression.Equals)
	}
	return false
}

func (Union) isNode() {}
