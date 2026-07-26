// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package modified applies tgen's own type corrections on top of the
// resolved specification: it introduces union and alias definitions the
// documentation does not name, and redirects matching field or method types
// to them.
package modified

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/resolved"
	"github.com/andreychh/tgen/model/pipeline/typed"
)

// Aliases is the table of tgen-introduced aliases, keyed by reference.
type Aliases = pipeline.Table[model.Reference, Alias]

// Specification is the database after tgen's own corrections are applied on
// top of the resolved stage: rules may introduce union or alias definitions
// the documentation does not name, and redirect matching field or method
// types to them.
type Specification struct {
	Objects  parsed.Objects
	Methods  resolved.Methods
	Fields   typed.Fields
	Unions   parsed.Unions
	Variants parsed.Variants
	Aliases  Aliases
	Release  parsed.Release
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
// modified one, applying tgen's own rules in order.
type Pass struct {
	spec resolved.Specification
}

// NewPass constructs a Pass over a resolved specification.
func NewPass(spec resolved.Specification) Pass {
	return Pass{spec: spec}
}

// Specification returns the modified specification, applying every rule tgen
// introduces, in order. No rule matches a type shape another rule produces, so
// the order is not significant; a rule added here must keep that true. It
// fails when any rule fails.
func (p Pass) Specification() (Specification, error) {
	spec := Specification{
		Objects:  p.spec.Objects,
		Methods:  p.spec.Methods,
		Fields:   p.spec.Fields,
		Unions:   p.spec.Unions,
		Variants: p.spec.Variants,
		Aliases:  pipeline.NewMapTable[model.Reference, Alias](),
		Release:  p.spec.Release,
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

// alreadyExists reports whether ref already names an object, method, union, or
// alias in spec. Methods count: a target renders each one as a type, so a
// method and a union sharing a reference collide there.
func alreadyExists(spec Specification, ref model.Reference) bool {
	if _, ok := spec.Objects.Lookup(ref); ok {
		return true
	}
	if _, ok := spec.Methods.Lookup(ref); ok {
		return true
	}
	if _, ok := spec.Unions.Lookup(ref); ok {
		return true
	}
	_, ok := spec.Aliases.Lookup(ref)
	return ok
}
