// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"slices"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/model/types"
)

// InputFile represents an Overlay that replaces InputFile-or-String and
// String-with-sending-files-link field types with InputFile.
type InputFile struct{}

func (o InputFile) Apply(field explicit.Field) explicit.Field {
	expr, err := field.Type().AsExpression()
	if err != nil {
		return field
	}
	links, err := field.Description().Links()
	if err != nil {
		return field
	}
	if expr.Equals(types.NewNamed("InputFile", types.KindObject)) ||
		expr.Equals(types.NewUnion(
			types.NewNamed("InputFile", types.KindObject),
			types.NewNamed("String", types.KindPrimitive),
		)) ||
		(expr.Equals(types.NewNamed("String", types.KindPrimitive)) &&
			slices.Contains(links, "#sending-files")) {
		return NewModified(field, types.NewNamed("InputFile", types.KindUnion))
	}
	return field
}
