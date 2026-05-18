// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/pkg/iters"
)

type DiscriminatedObject struct {
	inner ir.DiscriminatedObject
}

func NewDiscriminatedObject(v ir.DiscriminatedObject) DiscriminatedObject {
	return DiscriminatedObject{inner: v}
}

func (o DiscriminatedObject) Name() (string, error) {
	name, err := o.inner.Name()
	if err != nil {
		return "", err
	}
	return NewName(name).Value(), nil
}

func (o DiscriminatedObject) Doc() (string, error) {
	ref, err := o.inner.Reference()
	if err != nil {
		return "", err
	}
	doc, err := NewDefinitionDoc(ref, o.inner.Description()).Value()
	if err != nil {
		return "", err
	}
	return NewTypeGodoc(doc).Value(), nil
}

func (o DiscriminatedObject) Fields() DiscriminatedObjectFields {
	return NewDiscriminatedObjectFields(o.inner.Fields())
}

func (o DiscriminatedObject) HasInputFile() (bool, error) {
	return o.inner.HasInputFile()
}

func (o DiscriminatedObject) Unions() Unions {
	return Unions{inner: iters.NewMappedSeq(o.inner.Fields().Free(), NewField)}
}
