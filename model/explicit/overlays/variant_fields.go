// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// VariantFields represents the fields of a discriminated variant with overlay
// applied to free fields.
type VariantFields struct {
	inner   explicit.Fields
	overlay Overlay
}

// NewVariantFields constructs a VariantFields from parsed fields with the given
// overlay applied to free fields.
func NewVariantFields(f explicit.Fields, o Overlay) VariantFields {
	return VariantFields{inner: f, overlay: o}
}

func (f VariantFields) Free() iter.Seq[explicit.Field] {
	return iters.NewMappedSeq(f.inner.Free(), f.overlay.Apply)
}

func (f VariantFields) Discriminator() explicit.Discriminator {
	return f.inner.Discriminator()
}
