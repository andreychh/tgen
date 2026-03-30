// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import "github.com/andreychh/tgen/parsing"

// Sequential represents an Overlay that applies multiple overlays in order.
type Sequential struct {
	items []Overlay
}

// NewSequential constructs a Sequential from the given overlays.
func NewSequential(items ...Overlay) Sequential {
	return Sequential{items: items}
}

func (o Sequential) Apply(f parsing.Field) parsing.Field {
	for _, item := range o.items {
		f = item.Apply(f)
	}
	return f
}
