// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model/explicit"

// DiscriminatedVariant represents a single variant of a discriminated
// union for Go code generation.
type DiscriminatedVariant struct {
	inner explicit.DiscriminatedObject
}

// NewDiscriminatedVariant constructs a DiscriminatedVariant from a parsed variant object.
func NewDiscriminatedVariant(v explicit.DiscriminatedObject) DiscriminatedVariant {
	return DiscriminatedVariant{inner: v}
}

// Name returns the Go field name for this variant in the union struct.
func (v DiscriminatedVariant) Name() Name {
	return NewName(v.inner.Name())
}

func (v DiscriminatedVariant) DiscriminatorValue() DiscriminatorValue {
	return NewDiscriminatorValue(v.inner.Fields().Discriminator().Value())
}
