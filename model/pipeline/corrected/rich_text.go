// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

const (
	richTextRef         = model.Reference("richtext")
	richTextPlainRef    = model.Reference("richtextplain")
	richTextSequenceRef = model.Reference("richtextsequence")
)

// RichText is a [Rule] that introduces RichTextPlain and RichTextSequence —
// the two RichText variants the documentation writes as prose (a plain
// String, an Array of RichText) rather than as list items — and adds them as
// variants of the documented RichText union.
type RichText struct{}

// Apply implements [Rule]. It fails when richtextplain or richtextsequence
// already names something else in spec.
func (r RichText) Apply(spec Specification) (Specification, error) {
	definitions, err := r.definitions(spec.Definitions)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing rich text definitions: %w", err)
	}
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing rich text aliases: %w", err)
	}
	variants, err := r.variants(spec.Variants)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing rich text variants: %w", err)
	}
	return Specification{
		Definitions:    definitions,
		Methods:        spec.Methods,
		Fields:         spec.Fields,
		Discriminators: spec.Discriminators,
		Variants:       variants,
		Aliases:        aliases,
		Release:        spec.Release,
	}, nil
}

// definitions returns base holding what RichText names too: the two variants
// the documentation writes as prose rather than as list items. It fails when
// either reference is taken.
func (r RichText) definitions(base parsed.Definitions) (parsed.Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		richTextPlainRef,
		"RichTextPlain",
		model.DefinitionKindAlias,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"RichTextPlain represents the plain-text variant of a RichText value.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the rich text plain alias: %w", err)
	}
	err = out.Insert(
		richTextSequenceRef,
		"RichTextSequence",
		model.DefinitionKindAlias,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"RichTextSequence represents the nested-array variant of a RichText value.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the rich text sequence alias: %w", err)
	}
	return out, nil
}

// aliases returns the type each RichText variant stands for: a plain string,
// and an array of RichText itself.
func (r RichText) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](2)
	out.Insert(richTextPlainRef, Alias{
		Ref:  richTextPlainRef,
		Type: typeexpr.NewPrimitive(primitive.String),
	})
	out.Insert(richTextSequenceRef, Alias{
		Ref:  richTextSequenceRef,
		Type: typeexpr.NewArray(typeexpr.NewNamed(richTextRef)),
	})
	return out
}

// variants returns base listing RichTextPlain and RichTextSequence as variants
// of the documented RichText union too, after the ones it already lists. It
// fails when the union already lists either.
func (r RichText) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, richTextRef)
	err := out.Insert(
		richTextPlainRef,
		richTextSequenceRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listing rich text variants: %w", err)
	}
	return out, nil
}
