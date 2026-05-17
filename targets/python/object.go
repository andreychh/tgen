// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

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

func (o Object) Name() (ClassName, error) {
	name, err := o.inner.Name()
	if err != nil {
		return ClassName{}, err
	}
	return NewClassName(name), nil
}

func (o Object) Doc() (DocString, error) {
	ref, err := o.inner.Reference()
	if err != nil {
		return DocString{}, err
	}
	return NewClassDocString(ref, o.inner.Description()), nil
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
