// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model/ir"

// Discriminator represents the discriminator field of a variant for Go code generation.
type Discriminator struct {
	inner ir.Discriminator
}

// NewDiscriminator constructs a Discriminator from a parsed discriminator.
func NewDiscriminator(d ir.Discriminator) Discriminator {
	return Discriminator{inner: d}
}

// Key returns the discriminator field key (e.g. type from type="emoji").
func (d Discriminator) Key() (string, error) {
	key, err := d.inner.Key()
	if err != nil {
		return "", err
	}
	return string(key), nil
}

// Value returns the fixed discriminator value for this variant (e.g. "emoji").
func (d Discriminator) Value() (string, error) {
	val, err := d.inner.Value()
	if err != nil {
		return "", err
	}
	return "\"" + string(val) + "\"", nil
}
