// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/implicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type ImplicitDiscriminatedUnion struct {
	inner implicit.DiscriminatedUnion
}

func NewImplicitDiscriminatedUnion(u implicit.DiscriminatedUnion) ImplicitDiscriminatedUnion {
	return ImplicitDiscriminatedUnion{inner: u}
}

func (i ImplicitDiscriminatedUnion) Name() Name {
	return NewName(i.inner.Name())
}

func (i ImplicitDiscriminatedUnion) Doc() GoDoc {
	return NewGoDoc(i.inner.Description())
}

func (i ImplicitDiscriminatedUnion) DiscriminatorKey() Key {
	return NewKey(i.inner.DiscriminatorKey())
}

func (i ImplicitDiscriminatedUnion) Variants() iter.Seq[DiscriminatedVariant] {
	return iters.NewMappedSeq(
		i.inner.Variants(),
		func(v implicit.DiscriminatedVariant) DiscriminatedVariant {
			return NewImplicitDiscriminatedVariant(v)
		},
	)
}
