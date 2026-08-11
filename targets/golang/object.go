// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// Object represents the Go declaration of a documented object.
type Object struct {
	inner ir.Object
}

// NewObject creates an Object from the record of an object.
func NewObject(o ir.Object) Object {
	return Object{inner: o}
}

// Doc returns the doc comment of the declaration, closing with a link back to
// the section the object was read from.
func (o Object) Doc() string {
	return NewDefinitionDoc(o.inner.Ref, o.inner.Description, o.inner.Introduced).Value()
}

// Ref implements [Declaration].
func (o Object) Ref() string {
	return string(o.inner.Ref)
}

// Template implements [Declaration].
func (o Object) Template() string {
	return "object"
}

// Name returns the Go name the object declares.
func (o Object) Name() string {
	return NewName(o.inner.Name).Value()
}

// Fields returns the fields the object declares, in the order the documentation
// listed them.
func (o Object) Fields() []Field {
	return slices.NewMapped(o.inner.Fields, NewField)
}

// Unions returns the fields the object has to decode by hand, empty when
// encoding/json decodes it on its own.
func (o Object) Unions() []UnionField {
	return NewUnionFields(o.inner.Fields).Value()
}

// Files returns the fields the object has to hand a file over for, empty when
// it holds none and so rides encoding/json outbound too.
func (o Object) Files() []Attached {
	return slices.NewMapped(o.inner.Files, NewAttached)
}

// Rewrites reports whether the object has to rewrite itself into JSON because a
// union reaching a file admits it. An object holding a file of its own rewrites
// itself anyway; this is what obliges the ones holding none.
func (o Object) Rewrites() bool {
	return o.inner.Rewrites
}

// Direction returns which way the object travels, which is what decides the
// half of the codec it has to declare.
func (o Object) Direction() Direction {
	return NewDirection(o.inner.Direction)
}
