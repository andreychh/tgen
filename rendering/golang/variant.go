// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/parsing"

type Variant struct {
	inner parsing.Variant
}

func NewVariant(v parsing.Variant) Variant {
	return Variant{inner: v}
}

func (v Variant) Name() Name {
	return NewDefaultName(v.inner.Name())
}
