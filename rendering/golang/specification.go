// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/enrichment"
)

type Specification struct {
	inner enrichment.Specification
}

func NewSpecification(s enrichment.Specification) Specification {
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

func (s Specification) Methods() iter.Seq[Method] {
	return func(yield func(Method) bool) {
		for o := range s.inner.Methods() {
			if !yield(NewMethod(o)) {
				break
			}
		}
	}
}

func (s Specification) DiscriminatedUnions() iter.Seq[DiscriminatedUnion] {
	return func(yield func(DiscriminatedUnion) bool) {
		for u := range s.inner.DiscriminatedUnions() {
			if !yield(NewDiscriminatedUnion(u)) {
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

func (s Specification) ImplicitUnions() iter.Seq[ImplicitUnion] {
	return func(yield func(ImplicitUnion) bool) {
		for u := range s.inner.ImplicitUnions() {
			if !yield(NewImplicitUnion(u)) {
				break
			}
		}
	}
}

func (s Specification) Release() Release {
	return NewRelease(s.inner.Release())
}
