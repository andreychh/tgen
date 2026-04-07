// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/model/types"
)

// ChatID represents an Overlay that replaces Integer-or-String field types with
// ChatID.
type ChatID struct{}

func (o ChatID) Apply(field explicit.Field) explicit.Field {
	expr, err := field.Type().AsExpression()
	if err != nil {
		return field
	}
	if !expr.Equals(types.NewUnion(
		types.NewNamed("Integer", types.KindPrimitive),
		types.NewNamed("String", types.KindPrimitive),
	)) {
		return field
	}
	return NewModified(field, types.NewNamed("ChatID", types.KindUnion))
}
