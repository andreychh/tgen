// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/model/ir"
)

type Field struct {
	inner ir.Field
}

func NewField(f ir.Field) Field {
	return Field{inner: f}
}

func (f Field) Name() (string, error) {
	key, err := f.inner.Key()
	if err != nil {
		return "", err
	}
	return NewFieldName(key).Value(), nil
}

func (f Field) Annotation() (string, error) {
	typ, err := f.inner.Type()
	if err != nil {
		return "", err
	}
	opt, err := f.inner.Optionality()
	if err != nil {
		return "", err
	}
	return NewAnnotation(typ, opt).Value()
}

func (f Field) Doc() (string, error) {
	desc, err := f.inner.Description().Value()
	if err != nil {
		return "", err
	}
	return NewFieldDocString(desc).Value(), nil
}

func (f Field) Optional() (bool, error) {
	opt, err := f.inner.Optionality()
	if err != nil {
		return false, err
	}
	return bool(opt), nil
}

func (f Field) Key() (string, error) {
	key, err := f.inner.Key()
	if err != nil {
		return "", err
	}
	return string(key), nil
}

func (f Field) IsInputFile() (bool, error) {
	return f.inner.IsInputFile()
}

func (f Field) Part() (string, error) {
	typ, err := f.inner.Type()
	if err != nil {
		return "", err
	}
	part, err := NewType(typ).Part()
	if err != nil {
		return "", err
	}
	name, err := f.Name()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(part, "self."+name), nil
}
