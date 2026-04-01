// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model/implicit"

type ImplicitDiscriminatedVariant struct {
	inner implicit.DiscriminatedVariant
}

func NewImplicitDiscriminatedVariant(v implicit.DiscriminatedVariant) ImplicitDiscriminatedVariant {
	return ImplicitDiscriminatedVariant{inner: v}
}

func (i ImplicitDiscriminatedVariant) Name() Name {
	return NewName(i.inner.Name())
}

func (i ImplicitDiscriminatedVariant) DiscriminatorValue() DiscriminatorValue {
	return NewDiscriminatorValue(i.inner.DiscriminatorValue())
}
