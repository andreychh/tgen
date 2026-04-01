// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

import (
	"iter"
	"slices"

	"github.com/andreychh/tgen/model/implicit"
)

// Unions represents the set of implicitly defined union types.
type Unions struct {
	discriminated []implicit.DiscriminatedUnion
	structured    []implicit.StructuredUnion
}

// Discriminated returns an iterator over implicitly defined discriminated unions.
func (u Unions) Discriminated() iter.Seq[implicit.DiscriminatedUnion] {
	return slices.Values(u.discriminated)
}

// Structured returns an iterator over implicitly defined structured unions.
func (u Unions) Structured() iter.Seq[implicit.StructuredUnion] {
	return slices.Values(u.structured)
}
