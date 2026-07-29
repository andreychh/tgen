// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/resolved"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

const (
	maybeMessageRef = model.Reference("maybemessage")
	messageRef      = model.Reference("message")
	trueRef         = model.Reference("true")
)

// MaybeMessage is a [Rule] that introduces the MaybeMessage union — a method
// return value that is either an edited Message or True for inline
// messages — and redirects every method returning their union to it.
type MaybeMessage struct{}

// Apply implements [Rule]. It fails when maybemessage or true already names
// something else in spec.
func (r MaybeMessage) Apply(spec Specification) (Specification, error) {
	definitions, err := r.definitions(spec.Definitions)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing maybe message definitions: %w", err)
	}
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing maybe message aliases: %w", err)
	}
	variants, err := r.variants(spec.Variants)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing maybe message variants: %w", err)
	}
	methods, err := pipeline.NewMappedTable(spec.Methods, maybeMessageMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting maybe message methods: %w", err)
	}
	return Specification{
		Definitions:    definitions,
		Methods:        methods,
		Fields:         spec.Fields,
		Discriminators: spec.Discriminators,
		Variants:       variants,
		Aliases:        aliases,
		Release:        spec.Release,
	}, nil
}

// definitions returns base holding what MaybeMessage names too: the union
// itself and the True alias its variant stands for, which the documentation
// writes as a bare keyword rather than a named type. It fails when either
// reference is taken.
func (r MaybeMessage) definitions(base parsed.Definitions) (parsed.Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		maybeMessageRef,
		"MaybeMessage",
		model.DefinitionKindUnion,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"MaybeMessage represents a method return value that is either an "+
				"edited Message or True for inline messages.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the maybe message union: %w", err)
	}
	err = out.Insert(
		trueRef,
		"True",
		model.DefinitionKindAlias,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"True represents the boolean true value in Telegram API responses.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the true alias: %w", err)
	}
	return out, nil
}

// aliases returns the type the True alias stands for.
func (r MaybeMessage) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](1)
	out.Insert(trueRef, Alias{Ref: trueRef, Type: typeexpr.NewPrimitive(primitive.True)})
	return out
}

// variants returns base listing MaybeMessage's two variants, message and true.
// It fails when the union already lists either.
func (r MaybeMessage) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, maybeMessageRef)
	err := out.Insert(
		messageRef,
		trueRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listing maybe message variants: %w", err)
	}
	return out, nil
}

// maybeMessageMapping is a [pipeline.Mapping] that redirects a method returning
// the Message-or-True union to MaybeMessage; every other method rides through
// unchanged.
type maybeMessageMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (maybeMessageMapping) Apply(method resolved.Method) (resolved.Method, error) {
	if !method.Type.Equals(typeexpr.NewUnion(
		typeexpr.NewNamed(messageRef),
		typeexpr.NewPrimitive(primitive.True),
	)) {
		return method, nil
	}
	return resolved.Method{
		Ref:  method.Ref,
		Type: typeexpr.NewNamed(maybeMessageRef),
	}, nil
}
