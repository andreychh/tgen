// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedVariantFields groups the free fields and discriminator of a variant for Go code generation.
type DiscriminatedVariantFields struct {
	inner explicit.Fields
}

// NewVariantFields constructs a DiscriminatedVariantFields from parsed variant fields.
func NewVariantFields(f explicit.Fields) DiscriminatedVariantFields {
	return DiscriminatedVariantFields{inner: f}
}

// Free returns the fields that appear in the generated struct, excluding the discriminator.
func (f DiscriminatedVariantFields) Free() iter.Seq[Field] {
	return iters.NewMappedSeq(f.inner.Free(), NewField)
}

// Discriminator returns the discriminator field of this variant.
func (f DiscriminatedVariantFields) Discriminator() Discriminator {
	return NewDiscriminator(f.inner.Discriminator())
}
