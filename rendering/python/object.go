// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type Object struct {
	inner explicit.Object
}

func NewObject(o explicit.Object) Object {
	return Object{inner: o}
}

func (o Object) Name() ClassName {
	return NewClassName(o.inner.Name())
}

func (o Object) Doc() DocString {
	return NewClassDocString(o.inner.Reference(), o.inner.Description())
}

func (o Object) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(o.inner.Fields(), NewField)
}
