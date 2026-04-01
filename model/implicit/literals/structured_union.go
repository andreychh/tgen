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

// StructuredUnion represents an implicitly defined union type absent from the
// Telegram Bot API specification whose variants are distinguished by structure.
type StructuredUnion struct {
	name        string
	description string
	variants    []implicit.StructuredVariant
}

// NewStructuredUnion constructs a StructuredUnion from name, description, and
// variants.
func NewStructuredUnion(
	name, description string,
	variants []implicit.StructuredVariant,
) StructuredUnion {
	return StructuredUnion{
		name:        name,
		description: description,
		variants:    variants,
	}
}

func (u StructuredUnion) Name() model.Name {
	return literals.NewName(u.name)
}

func (u StructuredUnion) Description() model.Description {
	return literals.NewDescription(u.description, []string{})
}

func (u StructuredUnion) Variants() iter.Seq[implicit.StructuredVariant] {
	return slices.Values(u.variants)
}
