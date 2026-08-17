// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// DiscriminatedObject represents the Python declaration of an object one fixed
// value tells apart from the others its union stands for.
//
// The Go target writes that value into the JSON as it marshals, because a Go
// struct cannot declare a field that holds one constant. A Python field can, so
// the value is declared rather than written, and the declaration answers for
// both directions at once: pydantic reads the field to pick the variant on the
// way in, and dumps it on the way out. Declaring it is also the only way — a
// discriminated union refuses to build over a variant that leaves the key off.
type DiscriminatedObject struct {
	inner ir.DiscriminatedObject
}

// NewDiscriminatedObject creates a DiscriminatedObject from the record of a
// discriminated object.
func NewDiscriminatedObject(o ir.DiscriminatedObject) DiscriminatedObject {
	return DiscriminatedObject{inner: o}
}

// Doc returns the docstring of the declaration, closed by a link to the section
// the object was read from where the page gave it one.
func (o DiscriminatedObject) Doc() string {
	doc := NewDefinitionDoc(o.inner.Ref, o.inner.Description, o.inner.Introduced)
	return NewClassDocstring(doc.Passage()).Value()
}

// Ref implements [Declaration].
func (o DiscriminatedObject) Ref() string {
	return string(o.inner.Ref)
}

// Template implements [Declaration].
func (o DiscriminatedObject) Template() string {
	return "discriminated_object"
}

// Name implements [Declaration].
func (o DiscriminatedObject) Name() string {
	return NewClassName(o.inner.Name).Value()
}

// Discriminator returns the field holding the value the object is told apart
// by, which stands ahead of the rest since it is what names the variant.
func (o DiscriminatedObject) Discriminator() Discriminator {
	return NewDiscriminator(o.inner.Discriminator)
}

// Fields returns the fields the object declares beside the discriminator, the
// required ones before the optional, as the pipeline's exit ordered them.
func (o DiscriminatedObject) Fields() []Field {
	return slices.NewMapped(o.inner.Fields, NewField)
}

// Files returns the fields the object hands a file over for, empty when the
// object holds none.
func (o DiscriminatedObject) Files() []Attached {
	return slices.NewMapped(o.inner.Files, NewAttached)
}

// Rewrites reports whether the object has to answer for rewriting itself into
// JSON because a union reaching a file admits it. An object holding a file of
// its own rewrites itself anyway; this is what obliges the ones holding none.
func (o DiscriminatedObject) Rewrites() bool {
	return o.inner.Rewrites
}

// Direction returns which way the object travels, which is what decides whether
// the key a field is aliased by is the name a response is read by.
func (o DiscriminatedObject) Direction() Direction {
	return NewDirection(o.inner.Direction)
}

// Aliased reports whether any field of the object is read from a key its name
// stopped spelling. Such a class reads that key and nothing else, which holds
// only while the object travels inbound alone — hence the assertion the
// template pins on it.
func (o DiscriminatedObject) Aliased() bool {
	for _, field := range o.Fields() {
		if field.Aliased() {
			return true
		}
	}
	return false
}

// Discriminator represents the Python declaration of the field an object is
// told apart by: the one value its annotation admits, which is also the value
// it defaults to.
type Discriminator struct {
	inner ir.Discriminator
}

// NewDiscriminator creates a Discriminator from the record of a discriminator.
func NewDiscriminator(d ir.Discriminator) Discriminator {
	return Discriminator{inner: d}
}

// Name returns the Python field the discriminator declares.
func (d Discriminator) Name() string {
	return NewFieldName(d.inner.Key).Value()
}

// Annotation returns the annotation admitting the one value the object is told
// apart by.
func (d Discriminator) Annotation() string {
	return fmt.Sprintf("Literal[%q]", d.inner.Value)
}

// Assignment returns the value the field defaults to, which spares a caller
// naming what the class it is calling already says.
func (d Discriminator) Assignment() string {
	return fmt.Sprintf("%q", d.inner.Value)
}
