// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/pkg/iters"
)

type Method struct {
	inner ir.Method
}

func NewMethod(o ir.Method) Method {
	return Method{inner: o}
}

func (m Method) Name() (string, error) {
	name, err := m.inner.Name()
	if err != nil {
		return "", err
	}
	return NewName(name).Value(), nil
}

func (m Method) Doc() (string, error) {
	ref, err := m.inner.Reference()
	if err != nil {
		return "", err
	}
	doc, err := NewDefinitionDoc(ref, m.inner.Description()).Value()
	if err != nil {
		return "", err
	}
	return NewTypeGodoc(doc).Value(), nil
}

func (m Method) ReturnType() (Type, error) {
	typ, err := m.inner.ReturnType()
	if err != nil {
		return Type{}, err
	}
	return NewRequiredType(typ), nil
}

func (m Method) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(m.inner.Fields(), NewField)
}
