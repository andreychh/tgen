// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/literals"
	"github.com/andreychh/tgen/model/types"
)

// StructuredVariant represents a named variant of a StructuredUnion with a
// distinct field type.
type StructuredVariant struct {
	name string
	typ  string
	kind types.Kind
}

// NewStructuredVariant constructs a StructuredVariant from a field name and
// type name.
func NewStructuredVariant(name, typ string, kind types.Kind) StructuredVariant {
	return StructuredVariant{name: name, typ: typ, kind: kind}
}

func (v StructuredVariant) Name() model.Name {
	return literals.NewName(v.name)
}

func (v StructuredVariant) Type() model.Type {
	return literals.NewType(types.NewNamed(v.typ, v.kind))
}
