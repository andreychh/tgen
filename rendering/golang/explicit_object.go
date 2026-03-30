// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type ExplicitObject struct {
	inner explicit.Object
}

func NewExplicitObject(o explicit.Object) ExplicitObject {
	return ExplicitObject{inner: o}
}

func (o ExplicitObject) Name() Name {
	return NewName(o.inner.Name())
}

func (o ExplicitObject) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(o.inner.Reference(), o.inner.Description()))
}

func (o ExplicitObject) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(o.inner.Fields(), NewField)
}
