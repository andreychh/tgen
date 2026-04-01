// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/assembled"
	"github.com/andreychh/tgen/pkg/iters"
)

// Objects groups all object subtypes of the specification for Go code generation.
type Objects struct {
	inner assembled.Specification
}

// NewObjects constructs an Objects from an assembled specification.
func NewObjects(s assembled.Specification) Objects {
	return Objects{inner: s}
}

// Free returns all standalone objects from the explicit spec.
func (o Objects) Free() iter.Seq[Object] {
	return iters.NewMappedSeq(o.inner.Explicit().Objects(), NewObject)
}

// Discriminated returns all discriminated variant objects from the explicit spec.
func (o Objects) Discriminated() iter.Seq[DiscriminatedObject] {
	return func(yield func(DiscriminatedObject) bool) {
		for u := range o.inner.Explicit().Unions().Discriminated() {
			for v := range u.Variants() {
				if !yield(NewDiscriminatedObject(v)) {
					return
				}
			}
		}
	}
}
