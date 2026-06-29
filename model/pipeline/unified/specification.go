// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package unified merges object fields and method parameters into one table of
// uniformly shaped fields.
package unified

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// Fields is the table of unified fields, keyed by owner reference and field
// key.
type Fields = pipeline.Table[model.FieldKey, Field]

// Specification is the database after object fields and method parameters are
// merged into a single table of fields. The object, method, union, and variant
// tables and the release ride through from the parsed stage unchanged.
type Specification struct {
	Objects  parsed.Objects
	Methods  parsed.Methods
	Fields   Fields
	Unions   parsed.Unions
	Variants parsed.Variants
	Release  parsed.Release
}

// Pass is the unification stage: it rewrites a parsed specification into a
// unified one, merging object fields and method parameters into a single table.
type Pass struct {
	spec parsed.Specification
}

// NewPass constructs a Pass over a parsed specification.
func NewPass(spec parsed.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the unified specification, merging object fields and
// method parameters into one table. It fails when any field or parameter cannot
// be unified.
func (p Pass) Specification() (Specification, error) {
	fields, err := pipeline.NewMappedTable(p.spec.Fields, NewFieldMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("unifying fields: %w", err)
	}
	params, err := pipeline.NewMappedTable(p.spec.Params, NewParamMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("unifying parameters: %w", err)
	}
	merged, err := pipeline.NewMergedTable(fields, params).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("merging fields and parameters: %w", err)
	}
	return Specification{
		Objects:  p.spec.Objects,
		Methods:  p.spec.Methods,
		Fields:   merged,
		Unions:   p.spec.Unions,
		Variants: p.spec.Variants,
		Release:  p.spec.Release,
	}, nil
}
