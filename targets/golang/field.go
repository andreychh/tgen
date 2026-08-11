// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Field represents the Go declaration of a field an object owns or a parameter
// a method takes.
type Field struct {
	inner ir.Field
}

// NewField creates a Field from the record of a field.
func NewField(f ir.Field) Field {
	return Field{inner: f}
}

// Doc returns the doc comment of the declaration.
func (f Field) Doc() string {
	return NewFieldGodoc(f.inner.Description).Value()
}

// Name returns the Go name the field declares.
func (f Field) Name() string {
	return NewNameFromKey(f.inner.Key).Value()
}

// Type returns the Go type expression of the field.
func (f Field) Type() string {
	return NewType(f.inner.Type, f.inner.Optionality).Value()
}

// Tag returns the struct tag of the field.
func (f Field) Tag() string {
	return NewTag(f.inner.Key, f.inner.Optionality).Value()
}
