// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/model/types"
)

// ReplyMarkup represents an Overlay that replaces the four-variant reply markup
// union with ReplyMarkup.
type ReplyMarkup struct{}

func (o ReplyMarkup) Apply(field explicit.Field) explicit.Field {
	expr, err := field.Type().AsExpression()
	if err != nil {
		return field
	}
	if !expr.Equals(types.NewUnionType([]types.TypeExpression{
		types.NewNamedType("InlineKeyboardMarkup"),
		types.NewNamedType("ReplyKeyboardMarkup"),
		types.NewNamedType("ReplyKeyboardRemove"),
		types.NewNamedType("ForceReply"),
	})) {
		return field
	}
	return NewModified(field, types.NewNamedType("ReplyMarkup"))
}
