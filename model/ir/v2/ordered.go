// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"cmp"
	"slices"

	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// Ordered is the definitions table read as a sequence: what the table holds
// keyed by reference, given back in the order its source put it in.
type Ordered struct {
	definitions parsed.Definitions
}

// NewOrdered constructs an Ordered over a definitions table.
func NewOrdered(definitions parsed.Definitions) Ordered {
	return Ordered{definitions: definitions}
}

// Value returns every definition the table holds, ordered by the position its
// source gave it. A table iterates in no order, so ordering here is what makes
// generated output the same on every run.
func (o Ordered) Value() []parsed.Definition {
	out := make([]parsed.Definition, 0, o.definitions.Count())
	for _, definition := range o.definitions.All() {
		out = append(out, definition)
	}
	slices.SortFunc(out, func(a, b parsed.Definition) int {
		return cmp.Compare(a.Position, b.Position)
	})
	return out
}
