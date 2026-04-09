// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedObjectFields represents the fields of a discriminated variant
// with overlay applied to free fields.
type DiscriminatedObjectFields struct {
	inner   explicit.Fields
	overlay Overlay
}

// NewDiscriminatedObjectFields constructs a DiscriminatedObject from parsed
// fields with the given overlay applied to free fields.
func NewDiscriminatedObjectFields(f explicit.Fields, o Overlay) DiscriminatedObjectFields {
	return DiscriminatedObjectFields{inner: f, overlay: o}
}

func (f DiscriminatedObjectFields) Free() iter.Seq[explicit.Field] {
	return iters.NewMappedSeq(NewPrioritizedFields(f.inner.Free()), f.overlay.Apply)
}

func (f DiscriminatedObjectFields) Discriminator() explicit.Discriminator {
	return f.inner.Discriminator()
}
