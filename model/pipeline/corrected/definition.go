// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/prose"
)

// Definition is the record of anything a target renders, marked with where it
// came from: a heading of the documentation page, or a rule of this stage. Its
// reference, name, kind, position, and description carry over from the parsed
// stage. Introduced reports the latter, and only the former stands at a section
// a reference addresses: tgen chooses the reference of what it introduces, and
// a heading of that same name, where one stands, describes something else.
type Definition struct {
	Ref         model.Reference
	Name        model.Name
	Kind        model.DefinitionKind
	Position    model.Position
	Description prose.Passage
	Introduced  bool
}

// DefinitionMapping maps a parsed definition into a corrected one, marking it
// as the documentation's own. What this stage introduces is marked by the rule
// introducing it.
type DefinitionMapping struct{}

// NewDefinitionMapping constructs a DefinitionMapping.
func NewDefinitionMapping() DefinitionMapping {
	return DefinitionMapping{}
}

// Apply implements [pipeline.Mapping]. It never fails.
func (m DefinitionMapping) Apply(definition parsed.Definition) (Definition, error) {
	return Definition{
		Ref:         definition.Ref,
		Name:        definition.Name,
		Kind:        definition.Kind,
		Position:    definition.Position,
		Description: definition.Description,
		Introduced:  false,
	}, nil
}
