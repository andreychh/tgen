// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package typed resolves the prose of each field's type column into a type
// expression.
package typed

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/unified"
)

// Fields is the table of typed fields, keyed by owner reference and field key.
type Fields = pipeline.Table[model.FieldKey, Field]

// Specification is the database after every field's type prose is resolved into
// a type expression. The object, method, union, and variant tables and the
// release ride through from the unified stage unchanged.
type Specification struct {
	Objects  parsed.Objects
	Methods  parsed.Methods
	Fields   Fields
	Unions   parsed.Unions
	Variants parsed.Variants
	Release  parsed.Release
}

// Pass is the typing stage: it rewrites a unified specification into a typed
// one, resolving every field's type prose into a type expression.
type Pass struct {
	spec unified.Specification
}

// NewPass constructs a Pass over a unified specification.
func NewPass(spec unified.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the typed specification, resolving every field's type
// prose into a type expression. It fails when any field's type prose cannot be
// decoded.
func (p Pass) Specification() (Specification, error) {
	fields, err := pipeline.NewMappedTable(p.spec.Fields, NewFieldMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("typing fields: %w", err)
	}
	return Specification{
		Objects:  p.spec.Objects,
		Methods:  p.spec.Methods,
		Fields:   fields,
		Unions:   p.spec.Unions,
		Variants: p.spec.Variants,
		Release:  p.spec.Release,
	}, nil
}
