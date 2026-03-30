// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package implicit

import (
	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/parsing/explicit/literals"
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

func (v DiscriminatedVariant) Name() explicit.Name {
	return literals.NewName(v.name)
}

func (v DiscriminatedVariant) Type() explicit.Name {
	return literals.NewName(v.typ)
}
