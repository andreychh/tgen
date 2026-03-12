// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

type Specification struct {
	inner parsing.Specification
}

func NewSpecification(s parsing.Specification) Specification {
	return Specification{inner: s}
}

func (s Specification) Objects() iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for o := range s.inner.Objects() {
			if !yield(NewObject(o)) {
				break
			}
		}
	}
}

func (s Specification) Unions() iter.Seq[Union] {
	return func(yield func(Union) bool) {
		for u := range s.inner.Unions() {
			if !yield(NewUnion(u)) {
				break
			}
		}
	}
}

func (s Specification) Release() Release {
	return NewRelease(s.inner.Release())
}
