// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// InputMediaGroupRule represents an enrichment rule that replaces the input media array union with InputMediaGroup.
type InputMediaGroupRule struct{}

//nolint:ireturn // Field is the intentional public contract of Apply
func (r InputMediaGroupRule) Apply(field parsing.Field) parsing.Field {
	root, err := field.Type().Root()
	if err != nil {
		return field
	}
	if !root.Equal(parsing.NewArrayType(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("InputMediaAudio"),
		parsing.NewNamedType("InputMediaDocument"),
		parsing.NewNamedType("InputMediaPhoto"),
		parsing.NewNamedType("InputMediaVideo"),
	}))) {
		return field
	}
	return typedField{
		inner: field,
		tree: parsing.NewTypeTreeExpr(
			parsing.NewArrayType(
				parsing.NewNamedType("InputMediaGroup"),
			),
		),
	}
}
