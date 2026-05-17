// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
)

type Field struct {
	inner spec.Field
}

func NewField(f spec.Field) Field {
	return Field{inner: f}
}

func (f Field) Name() (FieldName, error) {
	key, err := f.inner.Key()
	if err != nil {
		return FieldName{}, err
	}
	return NewFieldName(model.Name(key)), nil
}

func (f Field) Annotation() (Annotation, error) {
	expr, err := f.inner.Type()
	if err != nil {
		return Annotation{}, err
	}
	opt, err := f.inner.Optionality()
	if err != nil {
		return Annotation{}, err
	}
	return NewAnnotation(expr, opt), nil
}

func (f Field) Doc() DocString {
	return NewFieldDocString(f.inner.Description())
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
	expr, err := f.inner.Type()
	if err != nil {
		return false, err
	}
	name, err := NewType(expr).name()
	if err != nil {
		return false, err
	}
	return name == "InputFile", nil
}

func (f Field) Part() (string, error) {
	expr, err := f.inner.Type()
	if err != nil {
		return "", err
	}
	part, err := NewType(expr).Part()
	if err != nil {
		return "", err
	}
	n, err := f.Name()
	if err != nil {
		return "", err
	}
	name, err := n.Value()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(part, "self."+name), nil
}
