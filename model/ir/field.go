// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/spec/overlays"
)

// Field represents a field of a Telegram Bot API object or method narrowed for
// code generation.
type Field struct {
	inner spec.Field
}

// NewField constructs a Field from a parsed field.
func NewField(f spec.Field) Field {
	return Field{inner: f}
}

func (f Field) Key() (model.Key, error) {
	return f.inner.Key()
}

func (f Field) Type() (Type, error) {
	expr, err := f.inner.Type()
	if err != nil {
		return Type{}, err
	}
	return NewType(expr), nil
}

func (f Field) Optionality() (model.Optionality, error) {
	return f.inner.Optionality()
}

func (f Field) Description() model.Description {
	return f.inner.Description()
}

func (f Field) IsInputFile() (bool, error) {
	t, err := f.Type()
	if err != nil {
		return false, err
	}
	name, err := t.Name()
	if err != nil {
		return false, err
	}
	return name == overlays.InputFileName, nil
}
