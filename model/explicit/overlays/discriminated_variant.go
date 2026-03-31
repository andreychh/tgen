// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
)

// DiscriminatedVariant represents an enriched discriminated variant with
// overlay applied to its free fields.
type DiscriminatedVariant struct {
	inner   explicit.DiscriminatedVariant
	overlay Overlay
}

// NewDiscriminatedVariant constructs a DiscriminatedVariant from a parsed
// variant with the given overlay applied to its free fields.
func NewDiscriminatedVariant(v explicit.DiscriminatedVariant, o Overlay) DiscriminatedVariant {
	return DiscriminatedVariant{inner: v, overlay: o}
}

func (v DiscriminatedVariant) Reference() model.Reference {
	return v.inner.Reference()
}

func (v DiscriminatedVariant) Name() model.Name {
	return v.inner.Name()
}

func (v DiscriminatedVariant) Description() model.Description {
	return v.inner.Description()
}

func (v DiscriminatedVariant) Fields() explicit.Fields {
	return NewVariantFields(v.inner.Fields(), v.overlay)
}
