// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/assembled"
	"github.com/andreychh/tgen/pkg/iters"
)

type Specification struct {
	inner assembled.Specification
}

func NewSpecification(s assembled.Specification) Specification {
	return Specification{inner: s}
}

func (s Specification) Objects() Objects {
	return NewObjects(s.inner)
}

func (s Specification) Methods() iter.Seq[Method] {
	return iters.NewMappedSeq(s.inner.Explicit().Methods(), NewMethod)
}

func (s Specification) Unions() Unions {
	return NewUnions(s.inner)
}

func (s Specification) Release() Release {
	return NewRelease(s.inner.Explicit().Release())
}
