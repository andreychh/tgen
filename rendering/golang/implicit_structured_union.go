// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/implicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type ImplicitStructuredUnion struct {
	inner implicit.StructuredUnion
}

func NewImplicitStructuredUnion(u implicit.StructuredUnion) ImplicitStructuredUnion {
	return ImplicitStructuredUnion{inner: u}
}

func (i ImplicitStructuredUnion) Name() Name {
	return NewName(i.inner.Name())
}

func (i ImplicitStructuredUnion) Doc() GoDoc {
	return NewGoDoc(i.inner.Description())
}

func (i ImplicitStructuredUnion) Variants() iter.Seq[StructuredVariant] {
	return iters.NewMappedSeq(
		i.inner.Variants(),
		func(v implicit.StructuredVariant) StructuredVariant {
			return NewImplicitStructuredVariant(v)
		},
	)
}
