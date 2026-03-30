// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package implicit

import (
	"iter"
	"slices"
)

// Unions represents the set of implicitly defined union types.
type Unions struct {
	discriminated []DiscriminatedUnion
}

// Discriminated returns an iterator over implicitly defined discriminated unions.
func (u Unions) Discriminated() iter.Seq[DiscriminatedUnion] {
	return slices.Values(u.discriminated)
}
