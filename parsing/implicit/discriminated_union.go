// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package implicit

import (
	"iter"
	"slices"

	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/parsing/explicit/literals"
)

// DiscriminatedUnion represents an implicitly defined union type absent from
// the Telegram Bot API specification.
type DiscriminatedUnion struct {
	name        string
	description string
	variants    []DiscriminatedVariant
}

// NewDiscriminatedUnion constructs a DiscriminatedUnion from name, description,
// and variants.
func NewDiscriminatedUnion(
	name, description string,
	variants []DiscriminatedVariant,
) DiscriminatedUnion {
	return DiscriminatedUnion{name: name, description: description, variants: variants}
}

func (u DiscriminatedUnion) Name() explicit.Name {
	return literals.NewName(u.name)
}

func (u DiscriminatedUnion) Description() explicit.Description {
	return literals.NewDescription(u.description, []string{})
}

func (u DiscriminatedUnion) Variants() iter.Seq[DiscriminatedVariant] {
	return slices.Values(u.variants)
}
