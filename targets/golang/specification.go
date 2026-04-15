// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type Specification struct {
	inner explicit.Specification
}

func NewSpecification(s explicit.Specification) Specification {
	return Specification{inner: s}
}

func (s Specification) Objects() iter.Seq[Object] {
	return iters.NewMappedSeq(s.inner.Objects(), NewObject)
}

func (s Specification) DiscriminatedObjects() iter.Seq[DiscriminatedObject] {
	return func(yield func(DiscriminatedObject) bool) {
		for u := range s.inner.DiscriminatedUnions() {
			for v := range u.Variants() {
				if !yield(NewDiscriminatedObject(v)) {
					return
				}
			}
		}
	}
}

func (s Specification) Methods() iter.Seq[Method] {
	return iters.NewMappedSeq(s.inner.Methods(), NewMethod)
}

func (s Specification) DiscriminatedUnions() iter.Seq[DiscriminatedUnion] {
	return iters.NewMappedSeq(s.inner.DiscriminatedUnions(), NewDiscriminatedUnion)
}

func (s Specification) Release() Release {
	return NewRelease(s.inner.Release())
}
