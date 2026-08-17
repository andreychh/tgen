// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// Object represents the Python declaration of a documented object.
type Object struct {
	inner ir.Object
}

// NewObject creates an Object from the record of an object.
func NewObject(o ir.Object) Object {
	return Object{inner: o}
}

// Doc returns the docstring of the declaration, closed by a link to the section
// the object was read from where the page gave it one.
func (o Object) Doc() string {
	doc := NewDefinitionDoc(o.inner.Ref, o.inner.Description, o.inner.Introduced)
	return NewClassDocstring(doc.Passage()).Value()
}

// Ref implements [Declaration].
func (o Object) Ref() string {
	return string(o.inner.Ref)
}

// Template implements [Declaration].
func (o Object) Template() string {
	return "object"
}

// Name implements [Declaration].
func (o Object) Name() string {
	return NewClassName(o.inner.Name).Value()
}

// Fields returns the fields the object declares, the required ones before the
// optional, as the pipeline's exit ordered them.
func (o Object) Fields() []Field {
	return slices.NewMapped(o.inner.Fields, NewField)
}

// Files returns the fields the object hands a file over for, empty when the
// object holds none.
func (o Object) Files() []Attached {
	return slices.NewMapped(o.inner.Files, NewAttached)
}

// Rewrites reports whether the object has to answer for rewriting itself into
// JSON because a union reaching a file admits it. An object holding a file of
// its own rewrites itself anyway; this is what obliges the ones holding none.
func (o Object) Rewrites() bool {
	return o.inner.Rewrites
}

// Direction returns which way the object travels, which is what decides whether
// the key a field is aliased by is the name a response is read by.
func (o Object) Direction() Direction {
	return NewDirection(o.inner.Direction)
}

// Aliased reports whether any field of the object is read from a key its name
// stopped spelling. Such a class reads that key and nothing else, which holds
// only while the object travels inbound alone — hence the assertion the
// template pins on it.
func (o Object) Aliased() bool {
	for _, field := range o.Fields() {
		if field.Aliased() {
			return true
		}
	}
	return false
}
