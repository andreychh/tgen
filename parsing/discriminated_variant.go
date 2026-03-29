// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import "github.com/andreychh/tgen/parsing/gq"

// GQDiscriminatedVariant represents a variant of a discriminated union parsed from the
// specification.
type GQDiscriminatedVariant struct {
	h4 gq.Selection
}

// NewGQDiscriminatedVariant constructs a GQDiscriminatedVariant from a h4 selection.
func NewGQDiscriminatedVariant(h4 gq.Selection) GQDiscriminatedVariant {
	return GQDiscriminatedVariant{h4: h4}
}

func (v GQDiscriminatedVariant) Name() Name {
	return NewGQName(v.h4)
}

func (v GQDiscriminatedVariant) Fields() Fields {
	return NewGQFields(v.h4)
}

func (v GQDiscriminatedVariant) Reference() Reference {
	return NewGQDefinitionReference(v.h4.Find("a.anchor"))
}

func (v GQDiscriminatedVariant) Description() Description {
	return NewGQDefinitionDescription(v.h4)
}
