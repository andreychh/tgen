// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// Object represents an enriched Telegram Bot API object definition.
type Object struct {
	inner   explicit.Object
	overlay Overlay
}

// NewObject constructs an Object from a parsed object with the given overlay
// applied to its fields.
func NewObject(o explicit.Object, overlay Overlay) Object {
	return Object{inner: o, overlay: overlay}
}

func (o Object) Reference() model.Reference {
	return o.inner.Reference()
}

func (o Object) Name() model.Name {
	return o.inner.Name()
}

func (o Object) Description() model.Description {
	return o.inner.Description()
}

func (o Object) Fields() iter.Seq[explicit.Field] {
	return iters.NewMappedSeq(NewPrioritizedFields(o.inner.Fields()), o.overlay.Apply)
}
