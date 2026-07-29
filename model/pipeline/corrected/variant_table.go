// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// VariantTable is the variants table a rule writes into, bound to the union it
// writes about: it places every variant inserted after all the variants that
// union already lists, so a rule states which types the union accepts and in
// which order, never a position. Its zero value is not usable; construct one
// with [NewVariantTable].
type VariantTable struct {
	base  parsed.Variants
	owner model.Reference
	added pipeline.MapTable[model.VariantKey, parsed.Variant]
}

// NewVariantTable constructs a VariantTable over the variants base already
// holds, writing about the union owner.
func NewVariantTable(base parsed.Variants, owner model.Reference) VariantTable {
	return VariantTable{
		base:  base,
		owner: owner,
		added: pipeline.NewMapTable[model.VariantKey, parsed.Variant](),
	}
}

// Insert admits the types the union accepts, each at the place after every
// variant the union already lists, in the order given. It fails when the union
// already lists one of them.
func (t VariantTable) Insert(refs ...model.Reference) error {
	for _, ref := range refs {
		key := model.VariantKey{Owner: t.owner, Ref: ref}
		if _, exists := t.Lookup(key); exists {
			return fmt.Errorf("%s already lists %s", t.owner, ref)
		}
		t.added.Insert(key, parsed.Variant{Ref: ref, Position: t.place()})
	}
	return nil
}

func (t VariantTable) Lookup(key model.VariantKey) (parsed.Variant, bool) {
	if variant, exists := t.added.Lookup(key); exists {
		return variant, true
	}
	return t.base.Lookup(key)
}

func (t VariantTable) Count() int {
	return t.base.Count() + t.added.Count()
}

func (t VariantTable) All() iter.Seq2[model.VariantKey, parsed.Variant] {
	return func(yield func(model.VariantKey, parsed.Variant) bool) {
		for key, variant := range t.base.All() {
			if !yield(key, variant) {
				return
			}
		}
		for key, variant := range t.added.All() {
			if !yield(key, variant) {
				return
			}
		}
	}
}

// place returns the position the next variant inserted takes: the one after
// every variant this table's union lists. Variants of other unions do not
// count, since each union numbers its own from zero.
func (t VariantTable) place() model.Position {
	next := model.Position(0)
	for key, variant := range t.All() {
		if key.Owner != t.owner {
			continue
		}
		next = max(next, variant.Position+1)
	}
	return next
}
