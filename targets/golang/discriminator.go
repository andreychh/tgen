// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model/spec"

// Discriminator represents the discriminator field of a variant for Go code generation.
type Discriminator struct {
	inner spec.Discriminator
}

// NewDiscriminator constructs a Discriminator from a parsed discriminator.
func NewDiscriminator(d spec.Discriminator) Discriminator {
	return Discriminator{inner: d}
}

// Key returns the discriminator field key (e.g. type from type="emoji").
func (d Discriminator) Key() (Key, error) {
	key, err := d.inner.Key()
	if err != nil {
		return Key{}, err
	}
	return NewKey(key), nil
}

// Value returns the fixed discriminator value for this variant (e.g. "emoji").
func (d Discriminator) Value() (DiscriminatorValue, error) {
	val, err := d.inner.Value()
	if err != nil {
		return DiscriminatorValue{}, err
	}
	return NewDiscriminatorValue(val), nil
}
