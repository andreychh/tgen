// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model/spec"

type DiscriminatedObject struct {
	inner spec.DiscriminatedObject
}

func NewDiscriminatedObject(v spec.DiscriminatedObject) DiscriminatedObject {
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
