// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// DiscriminatedObject represents the Go declaration of an object a union tells
// apart by a discriminator: the struct itself and the marshaller writing the
// fixed value it is told apart by. The value is written, never declared, so no
// variant can be built claiming to be another.
type DiscriminatedObject struct {
	inner ir.DiscriminatedObject
}

// NewDiscriminatedObject creates a DiscriminatedObject from the record of a
// discriminated object.
func NewDiscriminatedObject(o ir.DiscriminatedObject) DiscriminatedObject {
	return DiscriminatedObject{inner: o}
}

// Doc returns the doc comment of the declaration, closing with a link back to
// the section the object was read from.
func (o DiscriminatedObject) Doc() string {
	return NewDefinitionDoc(o.inner.Ref, o.inner.Description, o.inner.Introduced).Value()
}

// Ref implements [Declaration].
func (o DiscriminatedObject) Ref() string {
	return string(o.inner.Ref)
}

// Template implements [Declaration].
func (o DiscriminatedObject) Template() string {
	return "discriminated_object"
}

// Name returns the Go name the object declares.
func (o DiscriminatedObject) Name() string {
	return NewName(o.inner.Name).Value()
}

// Fields returns the fields the object declares, in the order the documentation
// listed them. The discriminating field is not among them.
func (o DiscriminatedObject) Fields() []Field {
	return slices.NewMapped(o.inner.Fields, NewField)
}

// Unions returns the fields the object has to decode by hand, empty when
// encoding/json decodes it on its own.
func (o DiscriminatedObject) Unions() []UnionField {
	return NewUnionFields(o.inner.Fields).Value()
}

// Files returns the fields the object has to hand a file over for, empty when
// it holds none and so rides its marshaller outbound.
func (o DiscriminatedObject) Files() []Attached {
	return slices.NewMapped(o.inner.Files, NewAttached)
}

// Rewrites reports whether the object has to rewrite itself into JSON because a
// union reaching a file admits it. An object holding a file of its own rewrites
// itself anyway; this is what obliges the ones holding none.
func (o DiscriminatedObject) Rewrites() bool {
	return o.inner.Rewrites
}

// Discriminator returns the field the object is marshalled with.
func (o DiscriminatedObject) Discriminator() Discriminator {
	return NewDiscriminator(o.inner.Discriminator)
}

// Direction returns which way the object travels, which is what decides the
// half of the codec it has to declare. An object no request carries never
// writes the value telling it apart, however plainly it is told apart.
func (o DiscriminatedObject) Direction() Direction {
	return NewDirection(o.inner.Direction)
}
