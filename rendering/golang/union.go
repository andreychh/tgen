// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

type Union struct {
	inner parsing.Union
}

func NewUnion(u parsing.Union) Union {
	return Union{inner: u}
}

func (u Union) Name() Name {
	return NewDefaultName(u.inner.Name())
}

func (u Union) Doc() Doc {
	return NewDoc(NewDefinitionDoc(u.inner.Ref(), u.inner.Description()))
}

func (u Union) Variants() iter.Seq[Variant] {
	return func(yield func(Variant) bool) {
		for v := range u.inner.Variants() {
			if !yield(NewVariant(v)) {
				break
			}
		}
	}
}
