// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// VariantFields groups the free fields and discriminator of a variant for Go code generation.
type VariantFields struct {
	inner parsing.VariantFields
}

// NewVariantFields constructs a VariantFields from parsed variant fields.
func NewVariantFields(f parsing.VariantFields) VariantFields {
	return VariantFields{inner: f}
}

// Free returns the fields that appear in the generated struct, excluding the discriminator.
func (f VariantFields) Free() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		for field := range f.inner.Free() {
			if !yield(NewField(field)) {
				break
			}
		}
	}
}

// Discriminator returns the discriminator field of this variant.
func (f VariantFields) Discriminator() Discriminator {
	return NewDiscriminator(f.inner.Discriminator())
}
