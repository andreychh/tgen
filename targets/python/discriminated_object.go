// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model/ir"

type DiscriminatedObject struct {
	inner ir.DiscriminatedObject
}

func NewDiscriminatedObject(v ir.DiscriminatedObject) DiscriminatedObject {
	return DiscriminatedObject{inner: v}
}

func (o DiscriminatedObject) Name() (ClassName, error) {
	name, err := o.inner.Name()
	if err != nil {
		return ClassName{}, err
	}
	return NewClassName(name), nil
}

func (o DiscriminatedObject) Doc() (DocString, error) {
	ref, err := o.inner.Reference()
	if err != nil {
		return DocString{}, err
	}
	return NewClassDocString(ref, o.inner.Description()), nil
}

func (o DiscriminatedObject) Fields() DiscriminatedObjectFields {
	return NewDiscriminatedObjectFields(o.inner.Fields())
}

func (o DiscriminatedObject) HasInputFile() (bool, error) {
	return o.inner.HasInputFile()
}
