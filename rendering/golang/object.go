// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/enrichment"
)

type Object struct {
	inner enrichment.Object
}

func NewObject(o enrichment.Object) Object {
	return Object{inner: o}
}

func (o Object) Name() Name {
	return NewDefaultName(o.inner.Name())
}

func (o Object) Doc() Doc {
	return NewDoc(NewDefinitionDoc(o.inner.Ref(), o.inner.Description()))
}

func (o Object) Fields() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		for f := range o.inner.Fields() {
			if !yield(NewField(f)) {
				break
			}
		}
	}
}
