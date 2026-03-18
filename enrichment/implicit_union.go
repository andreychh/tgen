// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// ImplicitUnion represents a union type introduced by tgen to name inline
// union types that the Telegram Bot API specification leaves anonymous.
type ImplicitUnion interface {
	Name() parsing.ObjectName
	Description() string
	Variants() iter.Seq[parsing.Variant]
}
