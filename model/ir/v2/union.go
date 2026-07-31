// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
)

// DiscriminatedUnion is the record of a union one key tells every variant of
// apart, joined with that key and the variants it stands for. It stands beside
// [Union], not on top of it: the two partition the unions the specification
// names, and no union is both. Carrier reports that a file is reached through
// the union, which every variant answers for by rewriting itself into JSON —
// something a target has to say where it declares what the union accepts.
type DiscriminatedUnion struct {
	Ref         model.Reference
	Name        model.Name
	Description prose.Passage
	Key         model.Key
	Variants    []DiscriminatedVariant
	Carrier     bool
}

func (DiscriminatedUnion) isDefinition() {}

// DiscriminatedVariant is the record of one variant of a discriminated union:
// the name of the type it stands for and the value it is told apart by. The
// value is unique among the variants of its union — that is what made the union
// discriminated.
type DiscriminatedVariant struct {
	Name  model.Name
	Value model.DiscriminatorValue
}

// Union is the record of a union no key tells the variants of apart, joined with
// the variants it stands for. Carrier reports the same as it does for
// [DiscriminatedUnion]. The tables say which types the union accepts and nothing
// about how to tell them apart, because nothing in the documentation does: one
// union is told apart by the shape of the JSON, another by trying its variants
// in turn, a third is only ever sent and never read. A target writes that by
// hand.
type Union struct {
	Ref         model.Reference
	Name        model.Name
	Description prose.Passage
	Variants    []Variant
	Carrier     bool
}

func (Union) isDefinition() {}

// Variant is the record of one variant of a union: the name of the type it
// stands for.
type Variant struct {
	Name model.Name
}
