// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type Method struct {
	inner explicit.Method
}

func NewMethod(o explicit.Method) Method {
	return Method{inner: o}
}

func (m Method) Name() (ClassName, error) {
	name, err := m.inner.Name()
	if err != nil {
		return ClassName{}, err
	}
	return NewClassName(name), nil
}

func (m Method) Doc() (DocString, error) {
	ref, err := m.inner.Reference()
	if err != nil {
		return DocString{}, err
	}
	return NewClassDocString(ref, m.inner.Description()), nil
}

func (m Method) ReturnType() (Type, error) {
	expr, err := m.inner.ReturnType()
	if err != nil {
		return Type{}, err
	}
	return NewType(expr), nil
}

func (m Method) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(m.inner.Fields(), NewField)
}
