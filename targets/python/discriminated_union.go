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

func (i DiscriminatedUnion) Name() (ClassName, error) {
	name, err := i.inner.Name()
	if err != nil {
		return ClassName{}, err
	}
	return NewClassName(name), nil
}

func (i DiscriminatedUnion) Doc() DocString {
	return NewDocString(i.inner.Description(), 0)
}

func (i DiscriminatedUnion) DiscriminatorKey() (Key, error) {
	key, err := i.inner.DiscriminatorKey()
	if err != nil {
		return Key{}, err
	}
	return NewKey(key), nil
}

func (i DiscriminatedUnion) Variants() Variants {
	return NewVariants(i.inner.Variants())
}
