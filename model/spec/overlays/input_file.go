// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"slices"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/types"
)

// InputFileName is the synthetic union name produced by the InputFile overlay.
const InputFileName = "InputFile"

// InputFile represents an Overlay that replaces InputFile-or-String and
// String-with-sending-files-link field types with InputFile.
type InputFile struct{}

func (o InputFile) Apply(field spec.Field) spec.Field {
	expr, err := field.Type()
	if err != nil {
		return field
	}
	links, err := field.Description().Links()
	if err != nil {
		return field
	}
	if expr.Equals(types.NewNamed(InputFileName, types.KindObject)) ||
		expr.Equals(types.NewUnion(
			types.NewNamed(InputFileName, types.KindObject),
			types.NewNamed("String", types.KindPrimitive),
		)) ||
		(expr.Equals(types.NewNamed("String", types.KindPrimitive)) &&
			slices.Contains(links, "#sending-files")) {
		return NewModified(field, types.NewNamed(InputFileName, types.KindUnion))
	}
	return field
}
