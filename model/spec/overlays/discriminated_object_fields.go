// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedObjectFields represents the fields of a discriminated variant
// with overlay applied to free fields.
type DiscriminatedObjectFields struct {
	inner   spec.Fields
	overlay Overlay
}

// NewDiscriminatedObjectFields constructs a DiscriminatedObject from parsed
// fields with the given overlay applied to free fields.
func NewDiscriminatedObjectFields(f spec.Fields, o Overlay) DiscriminatedObjectFields {
	return DiscriminatedObjectFields{inner: f, overlay: o}
}

func (f DiscriminatedObjectFields) Free() iter.Seq[spec.Field] {
	return iters.NewMappedSeq(NewPrioritizedFields(f.inner.Free()), f.overlay.Apply)
}

func (f DiscriminatedObjectFields) Discriminator() spec.Discriminator {
	return f.inner.Discriminator()
}
