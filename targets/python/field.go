// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/model/explicit"
)

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

func (f Field) Optional() (bool, error) {
	return f.inner.Optionality().AsBool()
}

func (f Field) Key() (string, error) {
	return f.inner.Key().AsString()
}

func (f Field) Part() (string, error) {
	part, err := NewType(f.inner.Type()).Part()
	if err != nil {
		return "", err
	}
	name, err := f.Name().AsString()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(part, "self."+name), nil
}
