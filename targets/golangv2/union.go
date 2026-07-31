// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// Union represents the Go declaration of a union nothing tells the variants of
// apart: the interface and the marker each variant carries. A decoder is not
// among them. Which variant a payload holds is not written anywhere the
// specification can be read from — one union is told apart by the shape of the
// JSON, another by trying its variants in turn, a third is only ever sent — so
// the decoder is written by hand, in the block claiming this reference.
type Union struct {
	inner ir.Union
}

// NewUnion creates a Union from the record of a union.
func NewUnion(u ir.Union) Union {
	return Union{inner: u}
}

// Doc returns the doc comment of the declaration, closing with a link back to
// the section the union was read from.
func (u Union) Doc() string {
	return NewDefinitionDoc(u.inner.Ref, u.inner.Description).Value()
}

// Ref implements [Declaration].
func (u Union) Ref() string {
	return string(u.inner.Ref)
}

// Template implements [Declaration].
func (u Union) Template() string {
	return "union"
}

// Name returns the Go name the union declares.
func (u Union) Name() string {
	return NewName(u.inner.Name).Value()
}

// Carrier reports whether a file is reached through the union, which puts the
// rewrite every variant answers for into the interface.
func (u Union) Carrier() bool {
	return u.inner.Carrier
}

// Variants returns the types the union stands for, in the order the
// documentation listed them.
func (u Union) Variants() []Variant {
	return slices.NewMapped(u.inner.Variants, NewVariant)
}

// Variant represents the Go declaration of one variant of a union: the marker
// tying the type it names to the interface.
type Variant struct {
	inner ir.Variant
}

// NewVariant creates a Variant from the record of a variant.
func NewVariant(v ir.Variant) Variant {
	return Variant{inner: v}
}

// Name returns the Go name of the type the variant names.
func (v Variant) Name() string {
	return NewName(v.inner.Name).Value()
}

// DiscriminatedUnion represents the Go declaration of a union one key tells
// every variant of apart: the interface, the marker each variant carries, and
// the decoder reading that key off the payload to pick one.
type DiscriminatedUnion struct {
	inner ir.DiscriminatedUnion
}

// NewDiscriminatedUnion creates a DiscriminatedUnion from the record of a
// discriminated union.
func NewDiscriminatedUnion(u ir.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u}
}

// Doc returns the doc comment of the declaration, closing with a link back to
// the section the union was read from.
func (u DiscriminatedUnion) Doc() string {
	return NewDefinitionDoc(u.inner.Ref, u.inner.Description).Value()
}

// Ref implements [Declaration].
func (u DiscriminatedUnion) Ref() string {
	return string(u.inner.Ref)
}

// Template implements [Declaration].
func (u DiscriminatedUnion) Template() string {
	return "discriminated_union"
}

// Name returns the Go name the union declares.
func (u DiscriminatedUnion) Name() string {
	return NewName(u.inner.Name).Value()
}

// Key returns the JSON key the decoder reads the variant off, quoted as a Go
// string literal.
func (u DiscriminatedUnion) Key() string {
	return fmt.Sprintf("%q", u.inner.Key)
}

// Carrier reports whether a file is reached through the union, which puts the
// rewrite every variant answers for into the interface.
func (u DiscriminatedUnion) Carrier() bool {
	return u.inner.Carrier
}

// Variants returns the types the union stands for, in the order the
// documentation listed them.
func (u DiscriminatedUnion) Variants() []DiscriminatedVariant {
	return slices.NewMapped(u.inner.Variants, NewDiscriminatedVariant)
}

// DiscriminatedVariant represents the Go declaration of one variant of a
// discriminated union: the marker tying the type it names to the interface, and
// the value the decoder picks it by.
type DiscriminatedVariant struct {
	inner ir.DiscriminatedVariant
}

// NewDiscriminatedVariant creates a DiscriminatedVariant from the record of a
// discriminated variant.
func NewDiscriminatedVariant(v ir.DiscriminatedVariant) DiscriminatedVariant {
	return DiscriminatedVariant{inner: v}
}

// Name returns the Go name of the type the variant names.
func (v DiscriminatedVariant) Name() string {
	return NewName(v.inner.Name).Value()
}

// Value returns the value the variant is told apart by, quoted as a Go string
// literal.
func (v DiscriminatedVariant) Value() string {
	return fmt.Sprintf("%q", v.inner.Value)
}
