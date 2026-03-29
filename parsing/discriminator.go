// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"github.com/andreychh/tgen/parsing/gq"
)

// GQDiscriminator represents the discriminator field of a variant, parsed from
// its table row.
type GQDiscriminator struct {
	tr gq.Selection
}

// NewGQDiscriminator constructs a GQDiscriminator from a table row td.
func NewGQDiscriminator(tr gq.Selection) GQDiscriminator {
	return GQDiscriminator{tr: tr}
}

// Key returns the field key of the discriminator field.
func (d GQDiscriminator) Key() Key {
	return NewGQKey(d.tr.Find("td").At(0))
}

// Value returns the discriminator value from the field description.
func (d GQDiscriminator) Value() DiscriminatorValue {
	return NewGQDiscriminatorValue(d.tr.Find("td").At(2))
}
