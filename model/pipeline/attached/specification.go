// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package attached marks every definition that can carry a file: the type a
// file is sent as, and everything holding one however deep. A target reads the
// mark to decide how a method sends its request, since a file travels as a
// multipart part while everything else travels as JSON.
package attached

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/classified"
	"github.com/andreychh/tgen/model/pipeline/corrected"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/pipeline/parsed"
)

// Files is the table of what each definition has to do with a file, keyed by
// reference. A definition the table does not hold carries no file.
type Files = pipeline.Table[model.Reference, File]

// Specification is the database after every definition that can carry a file
// is marked. The definition, method, field, discriminator, variant, and alias
// tables and the release ride through from the flattened stage unchanged.
type Specification struct {
	Definitions    corrected.Definitions
	Methods        flattened.Methods
	Fields         flattened.Fields
	Files          Files
	Discriminators classified.Discriminators
	Variants       parsed.Variants
	Aliases        flattened.Aliases
	Release        parsed.Release
}

// Pass is the attaching stage: it rewrites a flattened specification into an
// attached one, spreading the mark of a file from the type it is sent as to
// everything that can hold one. It reads the flat types a field and an alias
// were reduced to, and nothing a later stage decides, so it stands here.
type Pass struct {
	spec flattened.Specification
}

// NewPass constructs a Pass over a flattened specification.
func NewPass(spec flattened.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the attached specification, marking every definition
// that can carry a file. It fails when the specification names no type a file
// is sent as, which leaves the mark nothing to spread from and would quietly
// report that no method sends a file at all.
func (p Pass) Specification() (Specification, error) {
	if _, found := p.spec.Definitions.Lookup(corrected.InputFileRef); !found {
		return Specification{}, fmt.Errorf("%s names no definition", corrected.InputFileRef)
	}
	return Specification{
		Definitions: p.spec.Definitions,
		Methods:     p.spec.Methods,
		Fields:      p.spec.Fields,
		Files: NewFileTable(
			corrected.InputFileRef,
			NewFieldRule(p.spec.Fields),
			NewVariantRule(p.spec.Variants),
			NewAliasRule(p.spec.Aliases),
		).Table(),
		Discriminators: p.spec.Discriminators,
		Variants:       p.spec.Variants,
		Aliases:        p.spec.Aliases,
		Release:        p.spec.Release,
	}, nil
}
