// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

// Discriminator represents the discriminator field of a variant, parsed from
// its table row.
type Discriminator struct {
	tr gq.Selection
}

// NewDiscriminator constructs a Discriminator from a table row td.
func NewDiscriminator(tr gq.Selection) Discriminator {
	return Discriminator{tr: tr}
}

// Key returns the field key of the discriminator field.
func (d Discriminator) Key() explicit.Key {
	return NewKey(d.tr.Find("td").At(0))
}

// Value returns the discriminator value from the field description.
func (d Discriminator) Value() explicit.DiscriminatorValue {
	return NewDiscriminatorValue(d.tr.Find("td").At(2))
}
