// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

// DiscriminatedObjectFields provides the free fields and discriminator of a
// variant object.
type DiscriminatedObjectFields struct {
	root, h4 gq.Selection
}

// NewDiscriminatedObjectFields constructs a DiscriminatedObjectFields from an
// h4 selection.
func NewDiscriminatedObjectFields(root, h4 gq.Selection) DiscriminatedObjectFields {
	return DiscriminatedObjectFields{root: root, h4: h4}
}

// Free returns the fields of the variant that appear in the generated struct,
// excluding the discriminator field.
func (f DiscriminatedObjectFields) Free() iter.Seq[explicit.Field] {
	return func(yield func(explicit.Field) bool) {
		seq := f.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			FilterFunc(func(tr gq.Selection) bool {
				return NewFieldRow(tr).Kind() != FieldKindDiscriminator
			}).
			All()
		for tr := range seq {
			if !yield(NewObjectField(f.root, tr)) {
				break
			}
		}
	}
}

func (f DiscriminatedObjectFields) Discriminator() explicit.Discriminator {
	for tr := range f.h4.Until("h3, h4, hr").Find("table tbody tr").All() {
		if NewFieldRow(tr).Kind() == FieldKindDiscriminator {
			return NewDiscriminator(tr)
		}
	}
	panic(
		"DiscriminatedObjectFields: no discriminator field found; td is not a discriminated object",
	)
}
