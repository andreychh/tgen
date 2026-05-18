// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"iter"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// Fields represents the field set of a discriminated object: the free fields
// and the discriminator.
type Fields struct {
	inner spec.Fields
}

// NewFields constructs a Fields from a parsed field set.
func NewFields(f spec.Fields) Fields {
	return Fields{inner: f}
}

func (f Fields) Free() iter.Seq[Field] {
	return iters.NewMappedSeq(f.inner.Free(), NewField)
}

func (f Fields) Discriminator() Discriminator {
	return NewDiscriminator(f.inner.Discriminator())
}
