// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"
	"iter"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/model/types"
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
	result, err := m.inner.Result()
	if err != nil {
		return Type{}, err
	}
	switch result := result.(type) {
	case ir.Command:
		return NewRequiredType(ir.NewType(types.NewNamed("True", types.KindPrimitive))), nil
	case ir.Value:
		return NewRequiredType(result.Type()), nil
	default:
		return Type{}, fmt.Errorf("rendering method result: unexpected result type %T", result)
	}
}

func (m Method) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(m.inner.Fields(), NewField)
}
