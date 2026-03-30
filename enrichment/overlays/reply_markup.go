// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/types"
)

// ReplyMarkup represents an Overlay that replaces the four-variant reply markup
// union with ReplyMarkup.
type ReplyMarkup struct{}

func (o ReplyMarkup) Apply(f parsing.Field) parsing.Field {
	expr, err := f.Type().AsExpression()
	if err != nil {
		return f
	}
	if !expr.Equals(types.NewUnionType([]types.TypeExpression{
		types.NewNamedType("InlineKeyboardMarkup"),
		types.NewNamedType("ReplyKeyboardMarkup"),
		types.NewNamedType("ReplyKeyboardRemove"),
		types.NewNamedType("ForceReply"),
	})) {
		return f
	}
	return NewModified(f, types.NewNamedType("ReplyMarkup"))
}
