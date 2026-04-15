// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedUnion represents a discriminated union from the explicit spec for Go code generation.
type DiscriminatedUnion struct {
	inner explicit.DiscriminatedUnion
}

// NewDiscriminatedUnion constructs an DiscriminatedUnion from an explicit discriminated union.
func NewDiscriminatedUnion(u explicit.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u}
}

// Name returns the Go type name for this union.
func (u DiscriminatedUnion) Name() Name {
	return NewName(u.inner.Name())
}

// Doc returns the godoc comment for this union.
func (u DiscriminatedUnion) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(u.inner.Reference(), u.inner.Description()))
}

// DiscriminatorKey returns the JSON key used to discriminate variants.
func (u DiscriminatedUnion) DiscriminatorKey() Key {
	return NewKey(u.inner.DiscriminatorKey())
}

// Variants returns all variants of this union.
func (u DiscriminatedUnion) Variants() iter.Seq[DiscriminatedVariant] {
	return iters.NewMappedSeq(u.inner.Variants(), NewDiscriminatedVariant)
}
