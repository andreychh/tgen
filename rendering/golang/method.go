// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/enrichment"
)

type Method struct {
	inner enrichment.Method
}

func NewMethod(o enrichment.Method) Method {
	return Method{inner: o}
}

func (o Method) Name() Name {
	return NewDefaultName(o.inner.Name())
}

func (o Method) Doc() Doc {
	return NewDoc(NewDefinitionDoc(o.inner.Ref(), o.inner.Description()))
}

func (o Method) ReturnType() Type {
	return NewType(
		o.inner.Returns(),
		NewMethodResultName(NewDefaultName(o.inner.Name())),
	)
}

func (o Method) Fields() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		for f := range o.inner.Fields() {
			if !yield(NewField(f)) {
				break
			}
		}
	}
}
