// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// Union represents the Python declaration of a union nothing tells the variants
// of apart: the alternatives it stands for, and nothing besides.
//
// The Go target declares an interface here and leaves the decoder to a block
// written by hand, because which variant a payload holds is written nowhere the
// specification can be read from. Python needs no such block: pydantic is handed
// the alternatives and tells them apart by the shape of what arrives, so what
// the specification does say is the whole declaration.
type Union struct {
	inner ir.Union
}

// NewUnion creates a Union from the record of a union.
func NewUnion(u ir.Union) Union {
	return Union{inner: u}
}

// Doc returns the docstring of the declaration, closing with a link back to the
// section the union was read from.
func (u Union) Doc() string {
	doc := NewDefinitionDoc(u.inner.Ref, u.inner.Description, u.inner.Introduced)
	return NewStatementDocstring(doc.Passage()).Value()
}

// Ref implements [Declaration].
func (u Union) Ref() string {
	return string(u.inner.Ref)
}

// Template implements [Declaration].
func (u Union) Template() string {
	return "union"
}

// Name returns the Python name the union declares.
func (u Union) Name() string {
	return NewClassName(u.inner.Name).Value()
}

// Variants returns the names of the types the union stands for, in the order
// the documentation listed them.
func (u Union) Variants() []string {
	return slices.NewMapped(u.inner.Variants, variantName)
}

// DiscriminatedUnion represents the Python declaration of a union one key tells
// every variant of apart: the alternatives it stands for and the key telling
// them apart.
//
// The values that key holds are not written here, unlike in the Go target,
// whose decoder switches on them. Each variant declares its own value as the
// annotation of the discriminating field, and naming the key is all pydantic
// needs to read it off the payload and pick.
type DiscriminatedUnion struct {
	inner ir.DiscriminatedUnion
}

// NewDiscriminatedUnion creates a DiscriminatedUnion from the record of a
// discriminated union.
func NewDiscriminatedUnion(u ir.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u}
}

// Doc returns the docstring of the declaration, closing with a link back to the
// section the union was read from.
func (u DiscriminatedUnion) Doc() string {
	doc := NewDefinitionDoc(u.inner.Ref, u.inner.Description, u.inner.Introduced)
	return NewStatementDocstring(doc.Passage()).Value()
}

// Ref implements [Declaration].
func (u DiscriminatedUnion) Ref() string {
	return string(u.inner.Ref)
}

// Template implements [Declaration].
func (u DiscriminatedUnion) Template() string {
	return "discriminated_union"
}

// Name returns the Python name the union declares.
func (u DiscriminatedUnion) Name() string {
	return NewClassName(u.inner.Name).Value()
}

// Key returns the JSON key the variants are told apart by, quoted as a Python
// string literal.
func (u DiscriminatedUnion) Key() string {
	return fmt.Sprintf("%q", u.inner.Key)
}

// Variants returns the names of the types the union stands for, in the order
// the documentation listed them.
func (u DiscriminatedUnion) Variants() []string {
	return slices.NewMapped(u.inner.Variants, discriminatedVariantName)
}

// variantName returns the Python name of the type a variant names.
func variantName(v ir.Variant) string {
	return NewClassName(v.Name).Value()
}

// discriminatedVariantName returns the Python name of the type a discriminated
// variant names.
func discriminatedVariantName(v ir.DiscriminatedVariant) string {
	return NewClassName(v.Name).Value()
}
