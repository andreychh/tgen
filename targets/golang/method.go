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

func (m Method) Name() (Name, error) {
	name, err := m.inner.Name()
	if err != nil {
		return Name{}, err
	}
	return NewName(name), nil
}

func (m Method) Doc() (GoDoc, error) {
	ref, err := m.inner.Reference()
	if err != nil {
		return GoDoc{}, err
	}
	return NewGoDoc(NewDefinitionDoc(ref, m.inner.Description())), nil
}

func (m Method) ReturnType() (Type, error) {
	typ, err := m.inner.ReturnType()
	if err != nil {
		return Type{}, err
	}
	return NewType(typ, false), nil
}

func (m Method) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(m.inner.Fields(), NewField)
}
