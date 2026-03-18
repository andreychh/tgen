// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// ChatIdRule represents an enrichment rule that replaces Integer-or-String
// field types with ChatId.
type ChatIdRule struct{}

func (r ChatIdRule) Apply(f parsing.Field) parsing.Field {
	root, err := f.Type().Root()
	if err != nil {
		return f
	}
	if !root.Equal(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("Integer"),
		parsing.NewNamedType("String"),
	})) {
		return f
	}
	return typedField{
		inner: f,
		tree:  parsing.NewTypeTreeExpr(parsing.NewNamedType("ChatId")),
	}
}
