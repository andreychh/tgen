// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package classified decodes the fixed discriminator value each field's
// description may carry, moving it out of Fields into its own table.
package classified

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// Fields is the table of classified fields, keyed by owner reference and
// field key: every field whose description decodes a fixed discriminator
// value is excluded.
type Fields = pipeline.Table[model.FieldKey, Field]

// Discriminators is the table of fixed discriminator values decoded out of
// field descriptions, keyed by the reference of the field's owner.
type Discriminators = pipeline.Table[model.Reference, Discriminator]

// Specification is the database after every field's description is decoded
// for the fixed discriminator value it may carry: a field that decodes one
// moves from Fields into Discriminators. The object, method, parameter, union,
// and variant tables and the release ride through from the parsed stage
// unchanged.
type Specification struct {
	Objects        parsed.Objects
	Methods        parsed.Methods
	Fields         Fields
	Params         parsed.Params
	Discriminators Discriminators
	Unions         parsed.Unions
	Variants       parsed.Variants
	Release        parsed.Release
}

// Pass is the classification stage: it rewrites a parsed specification into a
// classified one, moving each field whose description decodes a fixed
// discriminator value out of Fields and into Discriminators.
type Pass struct {
	spec parsed.Specification
}

// NewPass constructs a Pass over a parsed specification.
func NewPass(spec parsed.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the classified specification, moving each field whose
// description decodes a fixed discriminator value out of Fields and into
// Discriminators.
func (p Pass) Specification() Specification {
	discriminators := NewDiscriminatorTable(p.spec.Fields).Apply()
	fields := pipeline.NewFilteredTable(p.spec.Fields, NewFieldFilter(discriminators)).Apply()
	return Specification{
		Objects:        p.spec.Objects,
		Methods:        p.spec.Methods,
		Fields:         fields,
		Params:         p.spec.Params,
		Discriminators: discriminators,
		Unions:         p.spec.Unions,
		Variants:       p.spec.Variants,
		Release:        p.spec.Release,
	}
}
