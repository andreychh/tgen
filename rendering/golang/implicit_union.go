// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/enrichment"
)

type ImplicitUnion struct {
	inner enrichment.ImplicitUnion
}

func NewImplicitUnion(u enrichment.ImplicitUnion) ImplicitUnion {
	return ImplicitUnion{inner: u}
}

func (u ImplicitUnion) Name() Name {
	return NewDefaultName(u.inner.Name())
}

func (u ImplicitUnion) Doc() Doc {
	return NewDoc(u.inner.Description())
}

func (u ImplicitUnion) Variants() iter.Seq[ImplicitVariant] {
	return func(yield func(ImplicitVariant) bool) {
		for v := range u.inner.Variants() {
			if !yield(NewImplicitVariant(v)) {
				break
			}
		}
	}
}
