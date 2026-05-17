// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
)

type Field struct {
	inner explicit.Field
}

func NewField(f explicit.Field) Field {
	return Field{inner: f}
}

func (f Field) Name() (Name, error) {
	key, err := f.inner.Key()
	if err != nil {
		return Name{}, err
	}
	return NewName(model.Name(key)), nil
}

func (f Field) Type() (Type, error) {
	expr, err := f.inner.Type()
	if err != nil {
		return nil, err
	}
	opt, err := f.inner.Optionality()
	if err != nil {
		return nil, err
	}
	return NewOptionalType(NewExprType(expr), opt), nil
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

func (f Field) Tag() (Tag, error) {
	key, err := f.inner.Key()
	if err != nil {
		return Tag{}, err
	}
	opt, err := f.inner.Optionality()
	if err != nil {
		return Tag{}, err
	}
	return NewTag(key, opt), nil
}

func (f Field) Doc() GoDoc {
	return NewGoDoc(f.inner.Description())
}

func (f Field) IsInputFile() (bool, error) {
	typ, err := f.Type()
	if err != nil {
		return false, err
	}
	name, err := typ.Name()
	if err != nil {
		return false, err
	}
	return name == "InputFile", nil
}

func (f Field) Part() (string, error) {
	typ, err := f.Type()
	if err != nil {
		return "", err
	}
	part, err := typ.Part()
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
	return fmt.Sprintf(part, "m."+name), nil
}
