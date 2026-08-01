// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package corrected applies tgen's own type corrections on top of the
// resolved specification: it introduces union and alias definitions the
// documentation does not name, and redirects matching field or method types
// to them.
//
// A definition tgen introduces has no place on the page, so it stands after
// every definition the page provides, in the order the rules run. Positions
// downstream of this stage therefore no longer index the page's headings.
//
// Every definition carries from here whether tgen introduced it. This is the
// last stage that can tell: downstream, a reference tgen chose is
// indistinguishable from one the page provides.
package corrected

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/classified"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/resolved"
	"github.com/andreychh/tgen/model/pipeline/typed"
)

// Definitions is the table of everything a target renders, keyed by reference,
// each marked with whether tgen introduced it.
type Definitions = pipeline.Table[model.Reference, Definition]

// Aliases is the table of what each alias has of its own, keyed by reference.
// Only a definition of alias kind has a record here.
type Aliases = pipeline.Table[model.Reference, Alias]

// Specification is the database after tgen's own corrections are applied on
// top of the resolved stage: rules may introduce union or alias definitions
// the documentation does not name, and redirect matching field or method
// types to them.
type Specification struct {
	Definitions    Definitions
	Methods        resolved.Methods
	Fields         typed.Fields
	Discriminators classified.Discriminators
	Variants       parsed.Variants
	Aliases        Aliases
	Release        parsed.Release
}

// Rule is a single tgen-introduced correction: given a specification, it
// returns the specification with its own definition introduced and any
// matching field or method type redirected to it.
type Rule interface {
	// Apply returns spec with this rule's correction applied. It fails when spec
	// already contains a definition the rule means to introduce.
	Apply(spec Specification) (Specification, error)
}

// Pass is the correction stage: it rewrites a resolved specification into a
// corrected one, applying tgen's own rules in order.
type Pass struct {
	spec resolved.Specification
}

// NewPass constructs a Pass over a resolved specification.
func NewPass(spec resolved.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the corrected specification: every definition the
// resolved stage provides marked as the documentation's own, then every rule
// tgen introduces applied in order. No rule matches a type shape another rule
// produces, so the order is not significant; a rule added here must keep that
// true. It fails when any rule fails.
func (p Pass) Specification() (Specification, error) {
	definitions, err := pipeline.NewMappedTable(p.spec.Definitions, NewDefinitionMapping()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("marking definitions: %w", err)
	}
	spec := Specification{
		Definitions:    definitions,
		Methods:        p.spec.Methods,
		Fields:         p.spec.Fields,
		Discriminators: p.spec.Discriminators,
		Variants:       p.spec.Variants,
		Aliases:        pipeline.NewMapTable[model.Reference, Alias](),
		Release:        p.spec.Release,
	}
	rules := []Rule{
		ChatID{},
		ReplyMarkup{},
		InputMediaGroup{},
		InputRichMedia{},
		InputFile{},
		MaybeMessage{},
		RichText{},
	}
	for _, rule := range rules {
		next, err := rule.Apply(spec)
		if err != nil {
			return Specification{}, fmt.Errorf("applying rule: %w", err)
		}
		spec = next
	}
	return spec, nil
}
