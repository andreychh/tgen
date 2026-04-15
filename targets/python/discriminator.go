// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model/explicit"

// Discriminator represents the discriminator field of a variant for Python code generation.
type Discriminator struct {
	inner explicit.Discriminator
}

// NewDiscriminator constructs a Discriminator from a parsed discriminator.
func NewDiscriminator(d explicit.Discriminator) Discriminator {
	return Discriminator{inner: d}
}

// Key returns the discriminator field key (e.g. type from type="emoji").
func (d Discriminator) Key() Key {
	return NewKey(d.inner.Key())
}

// Value returns the fixed discriminator value for this variant (e.g. "emoji").
func (d Discriminator) Value() DiscriminatorValue {
	return NewDiscriminatorValue(d.inner.Value())
}
