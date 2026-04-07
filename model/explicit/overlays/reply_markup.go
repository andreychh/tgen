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
	if !expr.Equals(types.NewUnion(
		types.NewNamed("InlineKeyboardMarkup", types.KindObject),
		types.NewNamed("ReplyKeyboardMarkup", types.KindObject),
		types.NewNamed("ReplyKeyboardRemove", types.KindObject),
		types.NewNamed("ForceReply", types.KindObject),
	)) {
		return field
	}
	return NewModified(field, types.NewNamed("ReplyMarkup", types.KindUnion))
}
