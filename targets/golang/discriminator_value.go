// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model"

// DiscriminatorValue represents the fixed discriminator string of a variant adapted for Go code generation.
type DiscriminatorValue struct {
	inner model.DiscriminatorValue
}

// NewDiscriminatorValue constructs a DiscriminatorValue from a parsed discriminator value.
func NewDiscriminatorValue(v model.DiscriminatorValue) DiscriminatorValue {
	return DiscriminatorValue{inner: v}
}

// AsString returns the discriminator value as a string.
func (v DiscriminatorValue) AsString() (string, error) {
	return v.inner.AsString()
}
