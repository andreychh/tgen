// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/prose"
)

// DefinitionTable is the definitions table a rule writes into: it places every
// definition inserted after all the definitions it already holds, so a rule
// states what it introduces and in which order, never a position. Its zero
// value is not usable; construct one with [NewDefinitionTable].
type DefinitionTable struct {
	base  parsed.Definitions
	added pipeline.MapTable[model.Reference, parsed.Definition]
}

// NewDefinitionTable constructs a DefinitionTable over the definitions base
// already holds.
func NewDefinitionTable(base parsed.Definitions) DefinitionTable {
	return DefinitionTable{
		base:  base,
		added: pipeline.NewMapTable[model.Reference, parsed.Definition](),
	}
}

// Insert admits a definition tgen introduces, at the place after every
// definition the table already holds. It fails when ref already names a
// definition, since a reference names at most one whatever its kind.
func (t DefinitionTable) Insert(
	ref model.Reference,
	name model.Name,
	kind model.DefinitionKind,
	description prose.Passage,
) error {
	if _, exists := t.Lookup(ref); exists {
		return fmt.Errorf("%s already names a definition", ref)
	}
	t.added.Insert(ref, parsed.Definition{
		Ref:         ref,
		Name:        name,
		Kind:        kind,
		Position:    t.place(),
		Description: description,
	})
	return nil
}

func (t DefinitionTable) Lookup(ref model.Reference) (parsed.Definition, bool) {
	if definition, exists := t.added.Lookup(ref); exists {
		return definition, true
	}
	return t.base.Lookup(ref)
}

func (t DefinitionTable) Count() int {
	return t.base.Count() + t.added.Count()
}

func (t DefinitionTable) All() iter.Seq2[model.Reference, parsed.Definition] {
	return func(yield func(model.Reference, parsed.Definition) bool) {
		for ref, definition := range t.base.All() {
			if !yield(ref, definition) {
				return
			}
		}
		for ref, definition := range t.added.All() {
			if !yield(ref, definition) {
				return
			}
		}
	}
}

// place returns the position the next definition inserted takes: the one after
// every definition the table holds.
func (t DefinitionTable) place() model.Position {
	next := model.Position(0)
	for _, definition := range t.All() {
		next = max(next, definition.Position+1)
	}
	return next
}
