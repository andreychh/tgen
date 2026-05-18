// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedObjectFields groups the free fields and discriminator of a variant for Go code generation.
type DiscriminatedObjectFields struct {
	inner ir.Fields
}

// NewDiscriminatedObjectFields constructs a DiscriminatedObjectFields from parsed variant fields.
func NewDiscriminatedObjectFields(f ir.Fields) DiscriminatedObjectFields {
	return DiscriminatedObjectFields{inner: f}
}

// Free returns the fields that appear in the generated struct, excluding the discriminator.
func (f DiscriminatedObjectFields) Free() iter.Seq[Field] {
	return iters.NewMappedSeq(f.inner.Free(), NewField)
}

// Discriminator returns the discriminator field of this variant.
func (f DiscriminatedObjectFields) Discriminator() Discriminator {
	return NewDiscriminator(f.inner.Discriminator())
}
