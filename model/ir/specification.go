// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package ir provides a narrowed view of the Telegram Bot API specification for
// code generation: field types are resolved into name, dimensionality, and kind.
package ir

import (
	"iter"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// Specification represents the Telegram Bot API specification narrowed for code
// generation.
type Specification struct {
	inner spec.Specification
}

// NewSpecification constructs a Specification from a parsed specification.
func NewSpecification(s spec.Specification) Specification {
	return Specification{inner: s}
}

func (s Specification) Objects() iter.Seq[Object] {
	return iters.NewMappedSeq(s.inner.Objects(), NewObject)
}

func (s Specification) Methods() iter.Seq[Method] {
	return iters.NewMappedSeq(s.inner.Methods(), NewMethod)
}

func (s Specification) DiscriminatedObjects() iter.Seq[DiscriminatedObject] {
	return iters.NewMappedSeq(s.inner.DiscriminatedObjects(), NewDiscriminatedObject)
}

func (s Specification) DiscriminatedUnions() iter.Seq[DiscriminatedUnion] {
	return iters.NewMappedSeq(s.inner.DiscriminatedUnions(), NewDiscriminatedUnion)
}

func (s Specification) Release() Release {
	return NewRelease(s.inner.Release())
}
