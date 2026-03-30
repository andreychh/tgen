// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/model/types"
)

// InputMediaGroup represents an Overlay that replaces the input media array
// union with InputMediaGroup.
type InputMediaGroup struct{}

func (o InputMediaGroup) Apply(field explicit.Field) explicit.Field {
	expr, err := field.Type().AsExpression()
	if err != nil {
		return field
	}
	if !expr.Equals(types.NewArrayType(types.NewUnionType([]types.TypeExpression{
		types.NewNamedType("InputMediaAudio"),
		types.NewNamedType("InputMediaDocument"),
		types.NewNamedType("InputMediaPhoto"),
		types.NewNamedType("InputMediaVideo"),
	}))) {
		return field
	}
	return NewModified(field, types.NewArrayType(types.NewNamedType("InputMediaGroup")))
}
