// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/model/ir"
)

type DiscriminatedUnion struct {
	inner ir.DiscriminatedUnion
}

func NewDiscriminatedUnion(u ir.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u}
}

func (i DiscriminatedUnion) Name() (string, error) {
	name, err := i.inner.Name()
	if err != nil {
		return "", err
	}
	return NewClassName(name).Value(), nil
}

func (i DiscriminatedUnion) Doc() (string, error) {
	ref, err := i.inner.Reference()
	if err != nil {
		return "", err
	}
	doc, err := NewDefinitionDoc(ref, i.inner.Description()).Value()
	if err != nil {
		return "", err
	}
	return NewClassDocString(doc).Value(), nil
}

func (i DiscriminatedUnion) DiscriminatorKey() (string, error) {
	key, err := i.inner.DiscriminatorKey()
	if err != nil {
		return "", err
	}
	return string(key), nil
}

func (i DiscriminatedUnion) Variants() (string, error) {
	return NewVariants(i.inner.Variants()).Value()
}
