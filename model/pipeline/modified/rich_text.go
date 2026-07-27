// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package modified

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/prose"
	typetree "github.com/andreychh/tgen/model/types/v2"
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
	for _, ref := range []model.Reference{richTextPlainRef, richTextSequenceRef} {
		if alreadyExists(spec, ref) {
			return Specification{}, fmt.Errorf("%s already names something else in spec", ref)
		}
	}
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing rich text aliases: %w", err)
	}
	variants, err := pipeline.NewMergedTable(spec.Variants, r.variants()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing rich text variants: %w", err)
	}
	return Specification{
		Objects:        spec.Objects,
		Methods:        spec.Methods,
		Fields:         spec.Fields,
		Discriminators: spec.Discriminators,
		Unions:         spec.Unions,
		Variants:       variants,
		Aliases:        aliases,
		Release:        spec.Release,
	}, nil
}

// aliases returns the table of RichText's two prose-only variants: a plain
// string, and an array of RichText itself.
func (r RichText) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](2)
	out.Insert(richTextPlainRef, Alias{
		Ref:  richTextPlainRef,
		Name: "RichTextPlain",
		Type: typetree.NewPrimitive(typetree.String),
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"RichTextPlain represents the plain-text variant of a RichText value.",
			prose.StylePlain,
		))),
	})
	out.Insert(richTextSequenceRef, Alias{
		Ref:  richTextSequenceRef,
		Name: "RichTextSequence",
		Type: typetree.NewArray(typetree.NewNamed(richTextRef)),
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"RichTextSequence represents the nested-array variant of a RichText value.",
			prose.StylePlain,
		))),
	})
	return out
}

// variants returns the table of RichTextPlain and RichTextSequence as
// variants of the documented RichText union.
func (r RichText) variants() parsed.Variants {
	out := pipeline.NewMapTableWithCapacity[model.VariantKey, parsed.Variant](2)
	out.Insert(
		model.VariantKey{Owner: richTextRef, Ref: richTextPlainRef},
		parsed.Variant{Ref: richTextPlainRef},
	)
	out.Insert(
		model.VariantKey{Owner: richTextRef, Ref: richTextSequenceRef},
		parsed.Variant{Ref: richTextSequenceRef},
	)
	return out
}
