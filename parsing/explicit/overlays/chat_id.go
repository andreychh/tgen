// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/parsing/types"
)

// ChatID represents an Overlay that replaces Integer-or-String field types with
// ChatID.
type ChatID struct{}

func (o ChatID) Apply(f explicit.Field) explicit.Field {
	expr, err := f.Type().AsExpression()
	if err != nil {
		return f
	}
	if !expr.Equals(types.NewUnionType([]types.TypeExpression{
		types.NewNamedType("Integer"),
		types.NewNamedType("String"),
	})) {
		return f
	}
	return NewModified(f, types.NewNamedType("ChatID"))
}
