// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

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
	return o.inner.HasInputFile()
}
