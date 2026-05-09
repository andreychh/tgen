// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type DiscriminatedObject struct {
	inner explicit.DiscriminatedObject
}

func NewDiscriminatedObject(v explicit.DiscriminatedObject) DiscriminatedObject {
	return DiscriminatedObject{inner: v}
}

func (o DiscriminatedObject) Name() Name {
	return NewName(o.inner.Name())
}

func (o DiscriminatedObject) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(o.inner.Reference(), o.inner.Description()))
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
