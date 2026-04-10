// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model/explicit"

type DiscriminatedObject struct {
	inner explicit.DiscriminatedObject
}

func NewDiscriminatedObject(v explicit.DiscriminatedObject) DiscriminatedObject {
	return DiscriminatedObject{inner: v}
}

func (o DiscriminatedObject) Name() ClassName {
	return NewClassName(o.inner.Name())
}

func (o DiscriminatedObject) Doc() DocString {
	return NewClassDocString(o.inner.Reference(), o.inner.Description())
}

func (o DiscriminatedObject) Fields() DiscriminatedObjectFields {
	return NewDiscriminatedObjectFields(o.inner.Fields())
}
