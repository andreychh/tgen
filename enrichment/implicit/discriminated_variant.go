// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package implicit

import (
	"github.com/andreychh/tgen/enrichment/literals"
	"github.com/andreychh/tgen/parsing"
)

// DiscriminatedVariant represents a named variant of a DiscriminatedUnion with
// a distinct field type.
type DiscriminatedVariant struct {
	name string
	typ  string
}

// NewDiscriminatedVariant constructs a DiscriminatedVariant from a field name
// and type name.
func NewDiscriminatedVariant(name, typ string) DiscriminatedVariant {
	return DiscriminatedVariant{name: name, typ: typ}
}

func (v DiscriminatedVariant) Name() parsing.Name {
	return literals.NewName(v.name)
}

func (v DiscriminatedVariant) Type() parsing.Name {
	return literals.NewName(v.typ)
}
