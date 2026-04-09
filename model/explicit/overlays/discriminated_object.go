// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
)

// DiscriminatedObject represents an enriched discriminated variant with
// overlay applied to its free fields.
type DiscriminatedObject struct {
	inner   explicit.DiscriminatedObject
	overlay Overlay
}

// NewDiscriminatedObject constructs a DiscriminatedObject from a parsed
// variant with the given overlay applied to its free fields.
func NewDiscriminatedObject(v explicit.DiscriminatedObject, o Overlay) DiscriminatedObject {
	return DiscriminatedObject{inner: v, overlay: o}
}

func (v DiscriminatedObject) Reference() model.Reference {
	return v.inner.Reference()
}

func (v DiscriminatedObject) Name() model.Name {
	return v.inner.Name()
}

func (v DiscriminatedObject) Description() model.Description {
	return v.inner.Description()
}

func (v DiscriminatedObject) Fields() explicit.Fields {
	return NewDiscriminatedObjectFields(v.inner.Fields(), v.overlay)
}
