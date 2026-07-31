// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"strings"

	ir "github.com/andreychh/tgen/model/ir/v2"
)

// UnionField represents a field its owner has to decode by hand. Nothing but a
// union needs one: encoding/json reaches into every other type on its own, but
// it cannot pick which variant an interface holds, so the field is shadowed by
// raw JSON and dispatched once the rest of the object is decoded. It is built
// only for a field whose type names a union, and so is never asked to shadow a
// field with nothing to shadow.
type UnionField struct {
	inner ir.Field
}

// NewUnionField creates a UnionField from the record of a field typed as a
// union.
func NewUnionField(f ir.Field) UnionField {
	return UnionField{inner: f}
}

// Name returns the Go name of the field, which the shadow declares too.
func (f UnionField) Name() string {
	return NewNameFromKey(f.inner.Key).Value()
}

// Tag returns the struct tag of the shadow. It never omits the field: a shadow
// is read and never written, so there is no encoding for omitempty to affect.
func (f UnionField) Tag() string {
	return NewTag(f.inner.Key, false).Value()
}

// Shadow returns the Go type holding the raw JSON of the field, enclosed in one
// slice per dimension so it shadows the field's own shape.
func (f UnionField) Shadow() string {
	return strings.Repeat("[]", int(f.inner.Type.Dimensionality())) + "json.RawMessage"
}

// Element returns the Go name of the union the raw JSON is dispatched into,
// without the slices enclosing it.
func (f UnionField) Element() string {
	return NewType(f.inner.Type, f.inner.Optionality).Name()
}

// Array reports whether the field holds a sequence, which is dispatched one
// element at a time.
func (f UnionField) Array() bool {
	return f.inner.Type.Dimensionality() > 0
}

// UnionFields is the view of the fields one owner has to decode by hand,
// narrowed from the fields it declares.
type UnionFields struct {
	fields []ir.Field
}

// NewUnionFields creates a UnionFields over the fields one owner declares.
func NewUnionFields(fields []ir.Field) UnionFields {
	return UnionFields{fields: fields}
}

// Value returns the fields typed as a union, in the order their owner declares
// them. It is empty when encoding/json decodes the owner on its own, which is
// what tells an owner needing a decoder from one that does not.
func (f UnionFields) Value() []UnionField {
	out := make([]UnionField, 0)
	for _, field := range f.fields {
		if !NewType(field.Type, field.Optionality).Union() {
			continue
		}
		out = append(out, NewUnionField(field))
	}
	return out
}
