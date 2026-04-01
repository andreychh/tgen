// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/model/implicit"
)

type ImplicitStructuredVariant struct {
	inner implicit.StructuredVariant
}

func NewImplicitStructuredVariant(inner implicit.StructuredVariant) ImplicitStructuredVariant {
	return ImplicitStructuredVariant{inner: inner}
}

func (i ImplicitStructuredVariant) Name() Name {
	return NewName(i.inner.Name())
}

func (i ImplicitStructuredVariant) Type() Type {
	return NewExprType(i.inner.Type())
}
