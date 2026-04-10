// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

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

func (s Specification) Release() Release {
	return NewRelease(s.inner.Release())
}
