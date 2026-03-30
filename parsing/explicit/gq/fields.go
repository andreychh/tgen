// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

// Fields provides the free fields and discriminator of a variant object.
type Fields struct {
	h4 gq.Selection
}

// NewFields constructs a Fields from an h4 selection.
func NewFields(h4 gq.Selection) Fields {
	return Fields{h4: h4}
}

// Free returns the fields of the variant that appear in the generated struct,
// excluding the discriminator field.
func (f Fields) Free() iter.Seq[explicit.Field] {
	return func(yield func(explicit.Field) bool) {
		seq := f.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			FilterFunc(func(tr gq.Selection) bool {
				return NewFieldRow(tr).Kind() != FieldKindDiscriminator
			}).
			All()
		for tr := range seq {
			if !yield(NewObjectField(tr)) {
				break
			}
		}
	}
}

func (f Fields) Discriminator() explicit.Discriminator {
	for tr := range f.h4.Until("h3, h4, hr").Find("table tbody tr").All() {
		if NewFieldRow(tr).Kind() == FieldKindDiscriminator {
			return NewDiscriminator(tr)
		}
	}
	panic("Fields: no discriminator field found; td is not a discriminated variant")
}
