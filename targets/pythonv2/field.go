// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Field represents the Python declaration of a field an object owns or a
// parameter a method takes.
type Field struct {
	inner ir.Field
}

// NewField creates a Field from the record of a field.
func NewField(f ir.Field) Field {
	return Field{inner: f}
}

// Doc returns the docstring of the declaration, which stands under it rather
// than over it, since a docstring is read from what precedes it.
func (f Field) Doc() string {
	return NewFieldDocstring(f.inner.Description).Value()
}

// Name returns the Python attribute the field declares.
func (f Field) Name() string {
	return NewAttribute(f.inner.Key).Value()
}

// Annotation returns the Python type annotation of the field.
func (f Field) Annotation() string {
	return NewAnnotation(f.inner.Type, f.inner.Optionality).Value()
}

// Assignment returns what the declaration stands equal to, empty where the
// annotation stands alone.
func (f Field) Assignment() string {
	return NewAssignment(f.inner.Key, f.inner.Optionality).Value()
}
