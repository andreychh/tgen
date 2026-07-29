// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package resolved decodes the return type of each method from its description
// prose into a type expression.
package resolved

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/classified"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/typed"
)

// Methods is the table of what each method has of its own, keyed by reference.
// Only a definition of method kind has a record here, so the table holds fewer
// records than there are definitions.
type Methods = pipeline.Table[model.Reference, Method]

// Specification is the database after every method's return type is decoded
// from its description prose into a type expression. The definition, field,
// discriminator, and variant tables and the release ride through from the
// typed stage unchanged.
type Specification struct {
	Definitions    parsed.Definitions
	Methods        Methods
	Fields         typed.Fields
	Discriminators classified.Discriminators
	Variants       parsed.Variants
	Release        parsed.Release
}

// Pass is the resolution stage: it rewrites a typed specification into a
// resolved one, decoding every method's return type from its description prose.
type Pass struct {
	spec typed.Specification
}

// NewPass constructs a Pass over a typed specification.
func NewPass(spec typed.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the resolved specification, decoding the return type of
// every definition of method kind from its description prose into a type
// expression. It fails when any method's return type prose cannot be decoded.
func (p Pass) Specification() (Specification, error) {
	definitions := pipeline.NewFilteredTable(
		p.spec.Definitions,
		parsed.NewKindFilter(model.DefinitionKindMethod),
	).Apply()
	methods, err := pipeline.NewMappedTable(definitions, NewMethodMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("resolving return types: %w", err)
	}
	return Specification{
		Definitions:    p.spec.Definitions,
		Methods:        methods,
		Fields:         p.spec.Fields,
		Discriminators: p.spec.Discriminators,
		Variants:       p.spec.Variants,
		Release:        p.spec.Release,
	}, nil
}
