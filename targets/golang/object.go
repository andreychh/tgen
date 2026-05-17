// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

type Object struct {
	inner spec.Object
}

func NewObject(o spec.Object) Object {
	return Object{inner: o}
}

func (o Object) Name() (Name, error) {
	name, err := o.inner.Name()
	if err != nil {
		return Name{}, err
	}
	return NewName(name), nil
}

func (o Object) Doc() (GoDoc, error) {
	ref, err := o.inner.Reference()
	if err != nil {
		return GoDoc{}, err
	}
	return NewGoDoc(NewDefinitionDoc(ref, o.inner.Description())), nil
}

func (o Object) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(o.inner.Fields(), NewField)
}

func (o Object) HasInputFile() (bool, error) {
	for f := range o.Fields() {
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

func (o Object) Unions() Unions {
	return Unions{inner: o.Fields()}
}
