// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/separated"
)

// Variants is the ordered view of the variants one union stands for, each named
// by the definition it addresses.
type Variants struct {
	db    separated.Specification
	owner model.Reference
}

// NewVariants constructs a Variants over the variants owner lists in the
// database.
func NewVariants(db separated.Specification, owner model.Reference) Variants {
	return Variants{db: db, owner: owner}
}

// Value returns the variants of the union, ordered by the position their source
// gave them. It fails when a variant names no definition.
func (v Variants) Value() ([]Variant, error) {
	records := v.records()
	out := make([]Variant, 0, len(records))
	for _, record := range records {
		definition, found := v.db.Definitions.Lookup(record.Ref)
		if !found {
			return nil, fmt.Errorf("variant %s names no definition", record.Ref)
		}
		out = append(out, Variant{Name: definition.Name})
	}
	return out, nil
}

// Discriminated returns the key the variants of the union are marked under and
// each variant with the value it is marked with, ordered by the position their
// source gave them. It returns no key and no variant when one of them is marked
// with nothing, when two are marked under different keys, or when two are marked
// with the same value: none of the three tells the variants apart, so the key
// means nothing and the union is read as a [Union]. It fails when a variant
// names no definition.
func (v Variants) Discriminated() (model.Key, []DiscriminatedVariant, error) {
	key := model.Key("")
	out := make([]DiscriminatedVariant, 0)
	for at, record := range v.records() {
		mark, marked := v.db.Discriminators.Lookup(record.Ref)
		if !marked || (at > 0 && mark.Key != key) || holds(out, mark.Value) {
			return "", nil, nil
		}
		definition, found := v.db.Definitions.Lookup(record.Ref)
		if !found {
			return "", nil, fmt.Errorf("variant %s names no definition", record.Ref)
		}
		key = mark.Key
		out = append(out, DiscriminatedVariant{Name: definition.Name, Value: mark.Value})
	}
	return key, out, nil
}

// records returns the raw variant records of the union, ordered by position.
func (v Variants) records() []parsed.Variant {
	out := make([]parsed.Variant, 0)
	for key, record := range v.db.Variants.All() {
		if key.Owner != v.owner {
			continue
		}
		out = append(out, record)
	}
	slices.SortFunc(out, func(a, b parsed.Variant) int {
		return cmp.Compare(a.Position, b.Position)
	})
	return out
}

// holds reports whether value already tells one of the variants apart, which no
// second variant of the same union may repeat.
func holds(variants []DiscriminatedVariant, value model.DiscriminatorValue) bool {
	return slices.ContainsFunc(variants, func(variant DiscriminatedVariant) bool {
		return variant.Value == value
	})
}
