// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/types"
)

// InputMediaGroup represents an Overlay that replaces the input media array
// union with InputMediaGroup.
type InputMediaGroup struct{}

func (o InputMediaGroup) Apply(f parsing.Field) parsing.Field {
	expr, err := f.Type().AsExpression()
	if err != nil {
		return f
	}
	if !expr.Equals(types.NewArrayType(types.NewUnionType([]types.TypeExpression{
		types.NewNamedType("InputMediaAudio"),
		types.NewNamedType("InputMediaDocument"),
		types.NewNamedType("InputMediaPhoto"),
		types.NewNamedType("InputMediaVideo"),
	}))) {
		return f
	}
	return NewModified(f, types.NewArrayType(types.NewNamedType("InputMediaGroup")))
}
