// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// ReplyMarkupRule represents an enrichment rule that replaces the four-variant reply markup union with ReplyMarkup.
type ReplyMarkupRule struct{}

//nolint:ireturn // Field is the intentional public contract of Apply
func (r ReplyMarkupRule) Apply(field parsing.Field) parsing.Field {
	root, err := field.Type().Root()
	if err != nil {
		return field
	}
	if !root.Equal(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("InlineKeyboardMarkup"),
		parsing.NewNamedType("ReplyKeyboardMarkup"),
		parsing.NewNamedType("ReplyKeyboardRemove"),
		parsing.NewNamedType("ForceReply"),
	})) {
		return field
	}
	return typedField{
		inner: field,
		tree:  parsing.NewTypeTreeExpr(parsing.NewNamedType("ReplyMarkup")),
	}
}
