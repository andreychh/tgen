// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// DiscriminatedUnion represents a discriminated union for Go code generation.
type DiscriminatedUnion struct {
	inner parsing.DiscriminatedUnion
}

// NewDiscriminatedUnion constructs a DiscriminatedUnion from a parsed
// discriminated union.
func NewDiscriminatedUnion(u parsing.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u}
}

// Name returns the Go type name for this union.
func (u DiscriminatedUnion) Name() Name {
	return NewDefaultName(u.inner.Name())
}

// Doc returns the godoc comment for this union.
func (u DiscriminatedUnion) Doc() Doc {
	return NewDoc(NewDefinitionDoc(u.inner.Ref(), u.inner.Description()))
}

// DiscriminatorKey returns the shared discriminator key for this union.
func (u DiscriminatedUnion) DiscriminatorKey() (string, error) {
	return u.inner.DiscriminatorKey().Value()
}

// Variants returns the variants of this union.
func (u DiscriminatedUnion) Variants() iter.Seq[VariantObject] {
	return func(yield func(VariantObject) bool) {
		for v := range u.inner.Variants() {
			if !yield(NewVariantObject(v)) {
				break
			}
		}
	}
}
