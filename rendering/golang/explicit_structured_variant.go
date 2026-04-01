// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model/explicit"

type ExplicitStructuredVariant struct {
	inner explicit.Object
}

func NewExplicitStructuredVariant(inner explicit.Object) ExplicitStructuredVariant {
	return ExplicitStructuredVariant{inner: inner}
}

func (v ExplicitStructuredVariant) Name() Name {
	return NewName(v.inner.Name())
}

func (v ExplicitStructuredVariant) Type() Type {
	return NewNamedType(v.inner.Name())
}
