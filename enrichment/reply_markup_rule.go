// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// ReplyMarkupRule represents an enrichment rule that replaces the four-variant reply markup union with ReplyMarkup.
type ReplyMarkupRule struct{}

func (r ReplyMarkupRule) Apply(f parsing.Field) parsing.Field {
	root, err := f.Type().Root()
	if err != nil {
		return f
	}
	if !root.Equal(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("InlineKeyboardMarkup"),
		parsing.NewNamedType("ReplyKeyboardMarkup"),
		parsing.NewNamedType("ReplyKeyboardRemove"),
		parsing.NewNamedType("ForceReply"),
	})) {
		return f
	}
	return typedField{
		inner: f,
		tree:  parsing.NewTypeTreeExpr(parsing.NewNamedType("ReplyMarkup")),
	}
}
