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
func (v DiscriminatedVariant) Name() (Name, error) {
	name, err := v.inner.Name()
	if err != nil {
		return Name{}, err
	}
	return NewName(name), nil
}

func (v DiscriminatedVariant) DiscriminatorValue() (DiscriminatorValue, error) {
	val, err := v.inner.Fields().Discriminator().Value()
	if err != nil {
		return DiscriminatorValue{}, err
	}
	return NewDiscriminatorValue(val), nil
}
