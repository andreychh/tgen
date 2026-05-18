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

func (o DiscriminatedObject) Name() (Name, error) {
	name, err := o.inner.Name()
	if err != nil {
		return Name{}, err
	}
	return NewName(name), nil
}

func (o DiscriminatedObject) Doc() (GoDoc, error) {
	ref, err := o.inner.Reference()
	if err != nil {
		return GoDoc{}, err
	}
	return NewGoDoc(NewDefinitionDoc(ref, o.inner.Description())), nil
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
