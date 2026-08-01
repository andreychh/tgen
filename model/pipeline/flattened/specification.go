// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package flattened reduces every type the pipeline carries to an atom and the
// number of dimensions over which that atom repeats, a form that cannot hold a
// union at all.
package flattened

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/classified"
	"github.com/andreychh/tgen/model/pipeline/corrected"
	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// Fields is the table of flattened fields, keyed by owner reference and field
// key.
type Fields = pipeline.Table[model.FieldKey, Field]

// Methods is the table of what each method has of its own, keyed by reference.
type Methods = pipeline.Table[model.Reference, Method]

// Aliases is the table of what each alias has of its own, keyed by reference.
type Aliases = pipeline.Table[model.Reference, Alias]

// Specification is the database after every type is reduced to a flat one. The
// definition, discriminator, and variant tables and the release ride through
// from the corrected stage unchanged.
type Specification struct {
	Definitions    corrected.Definitions
	Methods        Methods
	Fields         Fields
	Discriminators classified.Discriminators
	Variants       parsed.Variants
	Aliases        Aliases
	Release        parsed.Release
}

// Pass is the flattening stage: it rewrites a corrected specification into a
// flattened one, reducing every field, return, and alias type to a flat type.
type Pass struct {
	spec corrected.Specification
}

// NewPass constructs a Pass over a corrected specification.
func NewPass(spec corrected.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the flattened specification, reducing every type to an
// atom and a dimensionality. It fails when any type holds a union.
func (p Pass) Specification() (Specification, error) {
	fields, err := pipeline.NewMappedTable(p.spec.Fields, NewFieldMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("flattening fields: %w", err)
	}
	methods, err := pipeline.NewMappedTable(p.spec.Methods, NewMethodMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("flattening methods: %w", err)
	}
	aliases, err := pipeline.NewMappedTable(p.spec.Aliases, NewAliasMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("flattening aliases: %w", err)
	}
	return Specification{
		Definitions:    p.spec.Definitions,
		Methods:        methods,
		Fields:         fields,
		Discriminators: p.spec.Discriminators,
		Variants:       p.spec.Variants,
		Aliases:        aliases,
		Release:        p.spec.Release,
	}, nil
}
