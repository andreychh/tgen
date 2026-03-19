// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// ChatIdRule represents an enrichment rule that replaces Integer-or-String
// field types with ChatId.
type ChatIdRule struct{}

//nolint:ireturn // Field is the intentional public contract of Apply
func (r ChatIdRule) Apply(field parsing.Field) parsing.Field {
	root, err := field.Type().Root()
	if err != nil {
		return field
	}
	if !root.Equal(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("Integer"),
		parsing.NewNamedType("String"),
	})) {
		return field
	}
	return typedField{
		inner: field,
		tree:  parsing.NewTypeTreeExpr(parsing.NewNamedType("ChatId")),
	}
}
