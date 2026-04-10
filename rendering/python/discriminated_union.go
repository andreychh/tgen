// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/model/explicit"
)

type DiscriminatedUnion struct {
	inner explicit.DiscriminatedUnion
}

func NewDiscriminatedUnion(u explicit.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u}
}

func (i DiscriminatedUnion) Name() ClassName {
	return NewClassName(i.inner.Name())
}

func (i DiscriminatedUnion) Doc() DocString {
	return NewDocString(i.inner.Description(), 0)
}

func (i DiscriminatedUnion) DiscriminatorKey() Key {
	return NewKey(i.inner.DiscriminatorKey())
}

func (i DiscriminatedUnion) Variants() Variants {
	return NewVariants(i.inner.Variants())
}
