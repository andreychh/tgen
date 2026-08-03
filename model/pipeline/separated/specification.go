// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package separated splits each method's return into what the method signals:
// a command that only reports success, or a value that carries data.
package separated

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/attached"
	"github.com/andreychh/tgen/model/pipeline/classified"
	"github.com/andreychh/tgen/model/pipeline/corrected"
	"github.com/andreychh/tgen/model/pipeline/directed"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// Methods is the table of what each method has of its own, keyed by reference.
type Methods = pipeline.Table[model.Reference, Method]

// Specification is the database after every method's return is split into what
// the method signals. The definition, field, discriminator, variant, alias,
// file, and direction tables and the release ride through from the directed
// stage unchanged.
type Specification struct {
	Definitions    corrected.Definitions
	Methods        Methods
	Fields         flattened.Fields
	Files          attached.Files
	Directions     directed.Directions
	Discriminators classified.Discriminators
	Variants       parsed.Variants
	Aliases        flattened.Aliases
	Release        parsed.Release
}

// Pass is the separation stage: it rewrites a directed specification into a
// separated one, deciding for every method whether its return carries data or
// only reports success.
type Pass struct {
	spec directed.Specification
}

// NewPass constructs a Pass over a directed specification.
func NewPass(spec directed.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the separated specification, splitting every method's
// return into a command or a value. Classifying a return cannot fail, so the
// returned error is always nil.
func (p Pass) Specification() (Specification, error) {
	methods, err := pipeline.NewMappedTable(p.spec.Methods, NewMethodMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("separating methods: %w", err)
	}
	return Specification{
		Definitions:    p.spec.Definitions,
		Methods:        methods,
		Fields:         p.spec.Fields,
		Files:          p.spec.Files,
		Directions:     p.spec.Directions,
		Discriminators: p.spec.Discriminators,
		Variants:       p.spec.Variants,
		Aliases:        p.spec.Aliases,
		Release:        p.spec.Release,
	}, nil
}
