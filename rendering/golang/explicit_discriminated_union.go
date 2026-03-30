// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// ExplicitDiscriminatedUnion represents a discriminated union from the explicit spec for Go code generation.
type ExplicitDiscriminatedUnion struct {
	inner explicit.DiscriminatedUnion
}

// NewExplicitDiscriminatedUnion constructs an ExplicitDiscriminatedUnion from an explicit discriminated union.
func NewExplicitDiscriminatedUnion(u explicit.DiscriminatedUnion) ExplicitDiscriminatedUnion {
	return ExplicitDiscriminatedUnion{inner: u}
}

// Name returns the Go type name for this union.
func (u ExplicitDiscriminatedUnion) Name() Name {
	return NewName(u.inner.Name())
}

// Doc returns the godoc comment for this union.
func (u ExplicitDiscriminatedUnion) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(u.inner.Reference(), u.inner.Description()))
}

// DiscriminatorKey returns the JSON key used to discriminate variants.
func (u ExplicitDiscriminatedUnion) DiscriminatorKey() Key {
	return NewKey(u.inner.DiscriminatorKey())
}

// Variants returns all variants of this union.
func (u ExplicitDiscriminatedUnion) Variants() iter.Seq[DiscriminatedVariant] {
	return iters.NewMappedSeq(
		u.inner.Variants(),
		func(v explicit.DiscriminatedVariant) DiscriminatedVariant {
			return NewExplicitDiscriminatedVariant(v)
		},
	)
}
