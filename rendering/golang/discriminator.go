// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/parsing"

// Discriminator represents the discriminator field of a variant for Go code generation.
type Discriminator struct {
	inner parsing.Discriminator
}

// NewDiscriminator constructs a Discriminator from a parsed discriminator.
func NewDiscriminator(d parsing.Discriminator) Discriminator {
	return Discriminator{inner: d}
}

// Key returns the discriminator field key (e.g. type from type="emoji").
func (d Discriminator) Key() (string, error) {
	return d.inner.Key().Value()
}

// Value returns the fixed discriminator value for this variant (e.g. "emoji").
func (d Discriminator) Value() (string, error) {
	return d.inner.Value()
}
