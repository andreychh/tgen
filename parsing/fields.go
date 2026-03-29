// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

// GQFields provides the free fields and discriminator of a variant object.
type GQFields struct {
	h4 gq.Selection
}

// NewGQFields constructs a GQFields from an h4 td.
func NewGQFields(h4 gq.Selection) GQFields {
	return GQFields{h4: h4}
}

// Free returns the fields of the variant that appear in the generated struct,
// excluding the discriminator field.
func (f GQFields) Free() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		seq := f.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			FilterFunc(func(tr gq.Selection) bool {
				return NewGQFieldRow(tr).Kind() != FieldKindDiscriminator
			}).
			All()
		for tr := range seq {
			if !yield(NewGQObjectField(tr)) {
				break
			}
		}
	}
}

func (f GQFields) Discriminator() Discriminator {
	for tr := range f.h4.Until("h3, h4, hr").Find("table tbody tr").All() {
		if NewGQFieldRow(tr).Kind() == FieldKindDiscriminator {
			return NewGQDiscriminator(tr)
		}
	}
	panic("GQFields: no discriminator field found; td is not a discriminated variant")
}
