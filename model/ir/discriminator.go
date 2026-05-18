// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
)

// Discriminator represents the discriminator field of a discriminated object.
type Discriminator struct {
	inner spec.Discriminator
}

// NewDiscriminator constructs a Discriminator from a parsed discriminator.
func NewDiscriminator(d spec.Discriminator) Discriminator {
	return Discriminator{inner: d}
}

func (d Discriminator) Key() (model.Key, error) {
	return d.inner.Key()
}

func (d Discriminator) Value() (model.DiscriminatorValue, error) {
	return d.inner.Value()
}
