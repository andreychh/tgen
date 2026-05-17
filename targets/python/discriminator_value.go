// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model"

// DiscriminatorValue represents the fixed discriminator string of a variant adapted for Python code generation.
type DiscriminatorValue struct {
	inner model.DiscriminatorValue
}

// NewDiscriminatorValue constructs a DiscriminatorValue from a parsed discriminator value.
func NewDiscriminatorValue(v model.DiscriminatorValue) DiscriminatorValue {
	return DiscriminatorValue{inner: v}
}

func (v DiscriminatorValue) Value() (string, error) {
	return string(v.inner), nil
}
