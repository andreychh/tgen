// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// ImplicitVariant represents a named variant of an ImplicitUnion with a distinct field type.
type ImplicitVariant struct {
	name string
	typ  string
}

// NewImplicitVariant constructs an ImplicitVariant from a field name and type name.
func NewImplicitVariant(name, typ string) ImplicitVariant {
	return ImplicitVariant{name: name, typ: typ}
}

func NewTypeVariant(name string) ImplicitVariant {
	return NewImplicitVariant(name, name)
}

//nolint:ireturn // ObjectName is the intentional public contract of this method
func (v ImplicitVariant) Name() parsing.ObjectName {
	return staticObjectName{value: v.name}
}

//nolint:ireturn // ObjectName is the intentional public contract of this method
func (v ImplicitVariant) Type() parsing.ObjectName {
	return staticObjectName{value: v.typ}
}
