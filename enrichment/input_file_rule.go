// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"slices"

	"github.com/andreychh/tgen/parsing"
)

// InputFileRule represents an enrichment rule that replaces InputFile-or-String and String-with-sending-files-link field types with InputFile.
type InputFileRule struct{}

//nolint:ireturn // Field is the intentional public contract of Apply
func (r InputFileRule) Apply(field parsing.Field) parsing.Field {
	root, err := field.Type().Root()
	if err != nil {
		return field
	}
	if root.Equal(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("InputFile"),
		parsing.NewNamedType("String"),
	})) || root.Equal(parsing.NewNamedType("String")) &&
		slices.Contains(field.Description().Links(), "#sending-files") {
		return typedField{
			inner: field,
			tree:  parsing.NewTypeTreeExpr(parsing.NewNamedType("InputFile")),
		}
	}
	return field
}
