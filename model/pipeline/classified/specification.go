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
// value is excluded. Positions carry over from [parsed.Fields] unchanged, so
// an owner that lost a field has no field at position zero and its remaining
// fields are gapless only among themselves.
type Fields = pipeline.Table[model.FieldKey, Field]

// Discriminators is the table of fixed discriminator values decoded out of
// field descriptions, keyed by the reference of the field's owner.
type Discriminators = pipeline.Table[model.Reference, Discriminator]

// Specification is the database after every field's description is decoded
// for the fixed discriminator value it may carry: a field that decodes one
// moves from Fields into Discriminators. The definition, parameter, and
// variant tables and the release ride through from the parsed stage unchanged.
type Specification struct {
	Definitions    parsed.Definitions
	Fields         Fields
	Params         parsed.Params
	Discriminators Discriminators
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
// Discriminators. The error is always nil: nothing this pass does can fail,
// and it returns one only so that every pass in the chain is called the same
// way.
func (p Pass) Specification() (Specification, error) {
	discriminators := NewDiscriminatorTable(p.spec.Fields).Apply()
	fields := pipeline.NewFilteredTable(p.spec.Fields, NewFieldFilter(discriminators)).Apply()
	return Specification{
		Definitions:    p.spec.Definitions,
		Fields:         fields,
		Params:         p.spec.Params,
		Discriminators: discriminators,
		Variants:       p.spec.Variants,
		Release:        p.spec.Release,
	}, nil
}
