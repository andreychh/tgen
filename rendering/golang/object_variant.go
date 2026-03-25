// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/parsing"
)

// VariantObject represents a single variant of a discriminated union for Go code generation.
type VariantObject struct {
	inner parsing.VariantObject
}

// NewVariantObject constructs a VariantObject from a parsed variant object.
func NewVariantObject(v parsing.VariantObject) VariantObject {
	return VariantObject{inner: v}
}

// Name returns the Go field name for this variant in the union struct.
func (v VariantObject) Name() Name {
	return NewDefaultName(v.inner.Name())
}

// Type returns the Go struct type name for this variant.
func (v VariantObject) Type() Name {
	return NewDefaultName(v.inner.Name())
}

// Doc returns the godoc comment for this variant struct.
func (v VariantObject) Doc() Doc {
	return NewDoc(NewDefinitionDoc(v.inner.Ref(), v.inner.Description()))
}

// Fields returns the free fields and discriminator of this variant.
func (v VariantObject) Fields() VariantFields {
	return NewVariantFields(v.inner.Fields())
}
