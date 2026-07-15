// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/types"
)

// InputRichMedia represents an Overlay that replaces the input media union
// embedded in a rich message with InputRichMedia.
type InputRichMedia struct{}

func (o InputRichMedia) Apply(field spec.Field) spec.Field {
	expr, err := field.Type()
	if err != nil {
		return field
	}
	if !expr.Equals(types.NewUnion(
		types.NewNamed("InputMediaAnimation", types.KindObject),
		types.NewNamed("InputMediaAudio", types.KindObject),
		types.NewNamed("InputMediaPhoto", types.KindObject),
		types.NewNamed("InputMediaVideo", types.KindObject),
		types.NewNamed("InputMediaVoiceNote", types.KindObject),
	)) {
		return field
	}
	return NewModified(field, types.NewNamed("InputRichMedia", types.KindUnion))
}
