// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/assembled"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// Unions groups all union subtypes of the specification for Go code generation.
type Unions struct {
	inner assembled.Specification
}

// NewUnions constructs a Unions from an assembled specification.
func NewUnions(s assembled.Specification) Unions {
	return Unions{inner: s}
}

// Discriminated returns all discriminated unions from the explicit spec.
func (u Unions) Discriminated() iter.Seq[DiscriminatedUnion] {
	return iters.NewMappedSeq(
		u.inner.Explicit().Unions().Discriminated(),
		func(v explicit.DiscriminatedUnion) DiscriminatedUnion {
			return NewExplicitDiscriminatedUnion(v)
		},
	)
}

// Structured returns all structured unions from the explicit spec.
func (u Unions) Structured() iter.Seq[StructuredUnion] {
	return iters.NewMappedSeq(
		u.inner.Explicit().Unions().Structured(),
		func(v explicit.StructuredUnion) StructuredUnion {
			return NewExplicitStructuredUnion(v)
		},
	)
}
