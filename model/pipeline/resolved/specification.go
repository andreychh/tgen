// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package resolved decodes the return type of each method from its description
// prose into a [result.Result].
package resolved

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/typed"
)

// Methods is the table of resolved methods, keyed by reference.
type Methods = pipeline.Table[model.Reference, Method]

// Specification is the database after every method's return type is decoded
// from its description prose into a [result.Result]. The object, field, union,
// and variant tables and the release ride through from the typed stage
// unchanged.
type Specification struct {
	Objects  parsed.Objects
	Methods  Methods
	Fields   typed.Fields
	Unions   parsed.Unions
	Variants parsed.Variants
	Release  parsed.Release
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

// Specification returns the resolved specification, decoding every method's
// return type from its description prose into a [result.Result]. It fails when
// any method's return type prose cannot be decoded.
func (p Pass) Specification() (Specification, error) {
	methods, err := pipeline.NewMappedTable(p.spec.Methods, NewMethodMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("resolving methods: %w", err)
	}
	return Specification{
		Objects:  p.spec.Objects,
		Methods:  methods,
		Fields:   p.spec.Fields,
		Unions:   p.spec.Unions,
		Variants: p.spec.Variants,
		Release:  p.spec.Release,
	}, nil
}
