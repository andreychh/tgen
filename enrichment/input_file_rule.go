// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"slices"

	"github.com/andreychh/tgen/parsing"
)

// InputFileRule represents an enrichment rule that replaces InputFile-or-String and String-with-sending-files-link field types with InputFile.
type InputFileRule struct{}

func (r InputFileRule) Apply(f parsing.Field) parsing.Field {
	root, err := f.Type().Root()
	if err != nil {
		return f
	}
	if root.Equal(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("InputFile"),
		parsing.NewNamedType("String"),
	})) || root.Equal(parsing.NewNamedType("String")) &&
		slices.Contains(f.Description().Links(), "#sending-files") {
		return typedField{
			inner: f,
			tree:  parsing.NewTypeTreeExpr(parsing.NewNamedType("InputFile")),
		}
	}
	return f
}
