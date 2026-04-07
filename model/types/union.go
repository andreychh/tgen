// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"slices"
	"strings"

	slc "github.com/andreychh/tgen/pkg/slices"
)

type Union struct {
	variants []Expression
}

func NewUnion(variants ...Expression) Union {
	return Union{variants: variants}
}

func (u Union) Variants() []Expression {
	return u.variants
}

func (u Union) Equals(other Expression) bool {
	if other, ok := other.(Union); ok {
		return slices.EqualFunc(u.variants, other.variants, Expression.Equals)
	}
	return false
}

func (u Union) String() string {
	return strings.Join(slc.NewMapped(u.variants, Expression.String), " | ")
}

func (u Union) isNode() {}
