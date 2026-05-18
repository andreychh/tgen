// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedUnion represents a discriminated union for Go code generation.
type DiscriminatedUnion struct {
	inner ir.DiscriminatedUnion
}

// NewDiscriminatedUnion constructs a DiscriminatedUnion from a discriminated union.
func NewDiscriminatedUnion(u ir.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u}
}

// Name returns the Go type name for this union.
func (u DiscriminatedUnion) Name() (Name, error) {
	name, err := u.inner.Name()
	if err != nil {
		return Name{}, err
	}
	return NewName(name), nil
}

// Doc returns the godoc comment for this union.
func (u DiscriminatedUnion) Doc() (GoDoc, error) {
	ref, err := u.inner.Reference()
	if err != nil {
		return GoDoc{}, err
	}
	return NewGoDoc(NewDefinitionDoc(ref, u.inner.Description())), nil
}

// DiscriminatorKey returns the JSON key used to discriminate variants.
func (u DiscriminatedUnion) DiscriminatorKey() (Key, error) {
	key, err := u.inner.DiscriminatorKey()
	if err != nil {
		return Key{}, err
	}
	return NewKey(key), nil
}

// Variants returns all variants of this union.
func (u DiscriminatedUnion) Variants() iter.Seq[DiscriminatedVariant] {
	return iters.NewMappedSeq(u.inner.Variants(), NewDiscriminatedVariant)
}
