// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/pkg/iters"
)

// Object represents an enriched Telegram Bot API object definition.
type Object struct {
	inner   parsing.Object
	overlay Overlay
}

// NewObject constructs an Object from a parsed object with the given overlay
// applied to its fields.
func NewObject(o parsing.Object, overlay Overlay) Object {
	return Object{inner: o, overlay: overlay}
}

func (o Object) Reference() parsing.Reference {
	return o.inner.Reference()
}

func (o Object) Name() parsing.Name {
	return o.inner.Name()
}

func (o Object) Description() parsing.Description {
	return o.inner.Description()
}

func (o Object) Fields() iter.Seq[parsing.Field] {
	return iters.MapFunc(o.inner.Fields(), o.overlay.Apply)
}
