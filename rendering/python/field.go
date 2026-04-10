// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model/explicit"

type Field struct {
	inner explicit.Field
}

func NewField(f explicit.Field) Field {
	return Field{inner: f}
}

func (f Field) Name() FieldName {
	return NewFieldName(f.inner.Key())
}

func (f Field) Annotation() Annotation {
	return NewAnnotation(f.inner.Type(), f.inner.Optionality())
}

func (f Field) Doc() DocString {
	return NewFieldDocString(f.inner.Description())
}
