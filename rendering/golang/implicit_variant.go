// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/enrichment"

// ImplicitVariant represents a named variant of an ImplicitUnion for Go code generation.
type ImplicitVariant struct {
	inner enrichment.ImplicitVariant
}

// NewImplicitVariant constructs an ImplicitVariant from an enrichment variant.
func NewImplicitVariant(v enrichment.ImplicitVariant) ImplicitVariant {
	return ImplicitVariant{inner: v}
}

func (v ImplicitVariant) Name() Name {
	return NewDefaultName(v.inner.Name())
}

func (v ImplicitVariant) Type() VariantType {
	return NewVariantType(v.inner.Type())
}
