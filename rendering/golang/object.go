// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

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

func (o Object) Name() Name {
	return NewName(o.inner.Name())
}

func (o Object) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(o.inner.Reference(), o.inner.Description()))
}

func (o Object) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(o.inner.Fields(), NewField)
}
