// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

type DiscriminatedObject struct {
	inner spec.DiscriminatedObject
}

func NewDiscriminatedObject(v spec.DiscriminatedObject) DiscriminatedObject {
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
	for f := range o.Fields().Free() {
		ok, err := f.IsInputFile()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (o DiscriminatedObject) Unions() Unions {
	return Unions{inner: iters.NewMappedSeq(o.inner.Fields().Free(), NewField)}
}
