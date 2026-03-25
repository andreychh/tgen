// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

// GQVariantFields provides the free fields and discriminator of a variant object.
type GQVariantFields struct {
	selection gq.Selection
}

// NewGQVariantFields constructs a GQVariantFields from an h4 selection.
func NewGQVariantFields(h4 gq.Selection) GQVariantFields {
	return GQVariantFields{selection: h4}
}

// Free returns the fields of the variant that appear in the generated struct,
// excluding the discriminator field.
func (f GQVariantFields) Free() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		for tr := range f.selection.Until("h3, h4, hr").Find("table tbody tr").All() {
			if NewDefinitionRow(tr).Kind() == KindDiscriminatorField {
				continue
			}
			if !yield(NewObjectField(tr)) {
				break
			}
		}
	}
}

//nolint:ireturn // Discriminator is the intentional public contract of this method
func (f GQVariantFields) Discriminator() Discriminator {
	for tr := range f.selection.Until("h3, h4, hr").Find("table tbody tr").All() {
		if NewDefinitionRow(tr).Kind() == KindDiscriminatorField {
			return NewDiscriminator(tr)
		}
	}
	panic("GQVariantFields: no discriminator field found; selection is not a discriminated variant")
}
