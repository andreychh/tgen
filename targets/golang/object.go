// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/pkg/iters"
)

type Object struct {
	inner ir.Object
}

func NewObject(o ir.Object) Object {
	return Object{inner: o}
}

func (o Object) Name() (string, error) {
	name, err := o.inner.Name()
	if err != nil {
		return "", err
	}
	return NewName(name).Value(), nil
}

func (o Object) Doc() (string, error) {
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

func (o Object) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(o.inner.Fields(), NewField)
}

func (o Object) HasInputFile() (bool, error) {
	return o.inner.HasInputFile()
}

func (o Object) Unions() Unions {
	return Unions{inner: o.Fields()}
}
