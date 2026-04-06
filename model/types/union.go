// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

type Union struct {
	variants []Expression
}

func NewUnion(variants []Expression) Union {
	return Union{variants: variants}
}

func (u Union) Variants() []Expression {
	return u.variants
}

func (u Union) Equals(other Expression) bool {
	if other, ok := other.(Union); ok {
		if len(u.variants) != len(other.variants) {
			return false
		}
		for i, v := range u.variants {
			if !v.Equals(other.variants[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (u Union) isNode() {}
