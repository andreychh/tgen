// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

import (
	"iter"
	"slices"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/implicit"
	"github.com/andreychh/tgen/model/literals"
)

// DiscriminatedUnion represents an implicitly defined union type absent from
// the Telegram Bot API specification.
type DiscriminatedUnion struct {
	name        string
	description string
	key         string
	variants    []implicit.DiscriminatedVariant
}

// NewDiscriminatedUnion constructs a DiscriminatedUnion from name, description,
// and variants.
func NewDiscriminatedUnion(
	name, description, key string,
	variants []implicit.DiscriminatedVariant,
) DiscriminatedUnion {
	return DiscriminatedUnion{
		name:        name,
		description: description,
		key:         key,
		variants:    variants,
	}
}

func (u DiscriminatedUnion) Name() model.Name {
	return literals.NewName(u.name)
}

func (u DiscriminatedUnion) Description() model.Description {
	return literals.NewDescription(u.description, []string{})
}

func (u DiscriminatedUnion) DiscriminatorKey() model.Key {
	return literals.NewKey(u.key)
}

func (u DiscriminatedUnion) Variants() iter.Seq[implicit.DiscriminatedVariant] {
	return slices.Values(u.variants)
}
