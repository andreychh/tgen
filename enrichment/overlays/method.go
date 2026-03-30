// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/pkg/iters"
)

// Method represents an enriched Telegram Bot API method definition.
type Method struct {
	inner   parsing.Method
	overlay Overlay
}

// NewMethod constructs a Method from a parsed method with the given overlay
// applied to its fields.
func NewMethod(m parsing.Method, overlay Overlay) Method {
	return Method{inner: NewMaybeMessage(m), overlay: overlay}
}

func (m Method) Reference() parsing.Reference {
	return m.inner.Reference()
}

func (m Method) Name() parsing.Name {
	return m.inner.Name()
}

func (m Method) Description() parsing.Description {
	return m.inner.Description()
}

func (m Method) ReturnType() parsing.Type {
	return m.inner.ReturnType()
}

func (m Method) Fields() iter.Seq[parsing.Field] {
	return iters.MapFunc(m.inner.Fields(), m.overlay.Apply)
}
