// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model/explicit"

// ExplicitDiscriminatedVariant represents a single variant of a discriminated union for Go code generation.
type ExplicitDiscriminatedVariant struct {
	inner explicit.DiscriminatedVariant
}

// NewExplicitDiscriminatedVariant constructs a ExplicitDiscriminatedVariant from a parsed variant object.
func NewExplicitDiscriminatedVariant(v explicit.DiscriminatedVariant) ExplicitDiscriminatedVariant {
	return ExplicitDiscriminatedVariant{inner: v}
}

// Name returns the Go field name for this variant in the union struct.
func (v ExplicitDiscriminatedVariant) Name() Name {
	return NewName(v.inner.Name())
}

// Type returns the Go struct type name for this variant.
func (v ExplicitDiscriminatedVariant) Type() Name {
	return NewName(v.inner.Name())
}

// Doc returns the godoc comment for this variant struct.
func (v ExplicitDiscriminatedVariant) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(v.inner.Reference(), v.inner.Description()))
}

// Fields returns the free fields and discriminator of this variant.
func (v ExplicitDiscriminatedVariant) Fields() DiscriminatedVariantFields {
	return NewVariantFields(v.inner.Fields())
}
