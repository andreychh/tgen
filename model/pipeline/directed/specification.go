// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package directed decides which way every definition but a method travels:
// into a request, out of a response, or both. A target reads the direction to
// leave out an encoder or a decoder nothing could call, a definition no request
// carries never being written and one no response carries never being read.
package directed

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/attached"
	"github.com/andreychh/tgen/model/pipeline/classified"
	"github.com/andreychh/tgen/model/pipeline/corrected"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// Directions is the table of which way each definition travels, keyed by
// reference. A method is absent from it, being an operation rather than
// something carried: what travels is the parameters it takes and the value it
// hands back.
type Directions = pipeline.Table[model.Reference, model.Direction]

// Specification is the database after every definition is told which way it
// travels. The definition, method, field, file, discriminator, variant, and
// alias tables and the release ride through from the attached stage unchanged.
type Specification struct {
	Definitions    corrected.Definitions
	Methods        flattened.Methods
	Fields         flattened.Fields
	Files          attached.Files
	Directions     Directions
	Discriminators classified.Discriminators
	Variants       parsed.Variants
	Aliases        flattened.Aliases
	Release        parsed.Release
}

// Pass is the directing stage: it rewrites an attached specification into a
// directed one, spreading the way of travel from what a transport touches on
// its own to everything the specification reaches from there. It reads the flat
// types a field, a return, and an alias were reduced to, and nothing a later
// stage decides, so it stands here.
type Pass struct {
	spec attached.Specification
}

// NewPass constructs a Pass over an attached specification.
func NewPass(spec attached.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the directed specification, telling every definition
// but a method which way it travels. It fails when nothing reaches a
// definition, which leaves it neither written nor read and so declared for
// no reason at all.
func (p Pass) Specification() (Specification, error) {
	directions, err := p.directions(NewSpread(
		NewSeeding(p.spec.Methods, p.spec.Fields).Value(),
		NewFieldRule(p.spec.Fields),
		NewVariantRule(p.spec.Variants),
		NewAliasRule(p.spec.Aliases),
	).Value())
	if err != nil {
		return Specification{}, err
	}
	return Specification{
		Definitions:    p.spec.Definitions,
		Methods:        p.spec.Methods,
		Fields:         p.spec.Fields,
		Files:          p.spec.Files,
		Directions:     directions,
		Discriminators: p.spec.Discriminators,
		Variants:       p.spec.Variants,
		Aliases:        p.spec.Aliases,
		Release:        p.spec.Release,
	}, nil
}

// directions returns the way of travel of every definition the spread settled
// on, leaving out the methods. It fails when a definition that is no method was
// reached by nothing.
func (p Pass) directions(reaches Reaches) (Directions, error) {
	out := pipeline.NewMapTableWithCapacity[model.Reference, model.Direction](
		p.spec.Definitions.Count(),
	)
	for ref, definition := range p.spec.Definitions.All() {
		if definition.Kind == model.DefinitionKindMethod {
			continue
		}
		direction, err := reaches.Lookup(ref).Direction()
		if err != nil {
			return nil, fmt.Errorf("directing %s: %w", ref, err)
		}
		out.Insert(ref, direction)
	}
	return out, nil
}
