// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// Method represents an enriched Telegram Bot API method definition.
type Method struct {
	inner   explicit.Method
	overlay Overlay
}

// NewMethod constructs a Method from a parsed method with the given overlay
// applied to its fields.
func NewMethod(m explicit.Method, overlay Overlay) Method {
	return Method{inner: NewMaybeMessage(m), overlay: overlay}
}

func (m Method) Reference() explicit.Reference {
	return m.inner.Reference()
}

func (m Method) Name() explicit.Name {
	return m.inner.Name()
}

func (m Method) Description() explicit.Description {
	return m.inner.Description()
}

func (m Method) ReturnType() explicit.Type {
	return m.inner.ReturnType()
}

func (m Method) Fields() iter.Seq[explicit.Field] {
	return iters.NewMappedSeq(m.inner.Fields(), m.overlay.Apply)
}
