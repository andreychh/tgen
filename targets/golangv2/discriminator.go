// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/model/primitive"
)

// Discriminator represents the field a discriminated object is marshalled
// with. It is declared by no struct: the marshaller writes it into a shadow of
// the object, holding a constant no caller can set.
type Discriminator struct {
	inner ir.Discriminator
}

// NewDiscriminator creates a Discriminator from the record of a discriminator.
func NewDiscriminator(d ir.Discriminator) Discriminator {
	return Discriminator{inner: d}
}

// Name returns the Go name of the field.
func (d Discriminator) Name() string {
	return NewNameFromKey(d.inner.Key).Value()
}

// Type returns the Go type of the field. A discriminating value is a keyword
// the documentation writes as text, so Go declares it a string.
func (d Discriminator) Type() string {
	return builtin(primitive.String)
}

// Tag returns the struct tag of the field. The field is always written, so it
// is never omitted.
func (d Discriminator) Tag() string {
	return NewTag(d.inner.Key, false).Value()
}

// Value returns the constant the field holds, quoted as a Go string literal.
func (d Discriminator) Value() string {
	return fmt.Sprintf("%q", d.inner.Value)
}
