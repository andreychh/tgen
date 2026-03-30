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

func (o InputFile) Apply(f explicit.Field) explicit.Field {
	expr, err := f.Type().AsExpression()
	if err != nil {
		return f
	}
	links, err := f.Description().Links()
	if err != nil {
		return f
	}
	if expr.Equals(types.NewUnionType([]types.TypeExpression{
		types.NewNamedType("InputFile"),
		types.NewNamedType("String"),
	})) || expr.Equals(types.NewNamedType("String")) && slices.Contains(links, "#sending-files") {
		return NewModified(f, types.NewNamedType("InputFile"))
	}
	return f
}
