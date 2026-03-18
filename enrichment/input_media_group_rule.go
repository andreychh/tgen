// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// InputMediaGroupRule represents an enrichment rule that replaces the input media array union with InputMediaGroup.
type InputMediaGroupRule struct{}

func (r InputMediaGroupRule) Apply(f parsing.Field) parsing.Field {
	root, err := f.Type().Root()
	if err != nil {
		return f
	}
	if !root.Equal(parsing.NewArrayType(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("InputMediaAudio"),
		parsing.NewNamedType("InputMediaDocument"),
		parsing.NewNamedType("InputMediaPhoto"),
		parsing.NewNamedType("InputMediaVideo"),
	}))) {
		return f
	}
	return typedField{
		inner: f,
		tree:  parsing.NewTypeTreeExpr(parsing.NewArrayType(parsing.NewNamedType("InputMediaGroup"))),
	}
}
