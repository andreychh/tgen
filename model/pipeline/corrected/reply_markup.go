// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/typed"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

const (
	replyMarkupRef          = model.Reference("replymarkup")
	inlineKeyboardMarkupRef = model.Reference("inlinekeyboardmarkup")
	replyKeyboardMarkupRef  = model.Reference("replykeyboardmarkup")
	replyKeyboardRemoveRef  = model.Reference("replykeyboardremove")
	forceReplyRef           = model.Reference("forcereply")
)

// ReplyMarkup is a [Rule] that introduces the ReplyMarkup union — the four
// kinds of reply markup a message may carry — and redirects every field
// typed as their union to it.
type ReplyMarkup struct{}

// Apply implements [Rule]. It fails when replymarkup already names something
// else in spec.
func (r ReplyMarkup) Apply(spec Specification) (Specification, error) {
	definitions, err := r.definitions(spec.Definitions)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing reply markup definitions: %w", err)
	}
	variants, err := r.variants(spec.Variants)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing reply markup variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, replyMarkupMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting reply markup fields: %w", err)
	}
	return Specification{
		Definitions:    definitions,
		Methods:        spec.Methods,
		Fields:         fields,
		Discriminators: spec.Discriminators,
		Variants:       variants,
		Aliases:        spec.Aliases,
		Release:        spec.Release,
	}, nil
}

// definitions returns base holding what ReplyMarkup names too: the union
// itself. Its variants are documented objects and name themselves. It fails
// when the reference is taken.
func (r ReplyMarkup) definitions(base parsed.Definitions) (parsed.Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		replyMarkupRef,
		"ReplyMarkup",
		model.DefinitionKindUnion,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"ReplyMarkup represents a reply markup attached to a message.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the reply markup union: %w", err)
	}
	return out, nil
}

// variants returns base listing ReplyMarkup's four variants, each already a
// documented object. It fails when the union already lists one of them.
func (r ReplyMarkup) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, replyMarkupRef)
	err := out.Insert(
		inlineKeyboardMarkupRef,
		replyKeyboardMarkupRef,
		replyKeyboardRemoveRef,
		forceReplyRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listing reply markup variants: %w", err)
	}
	return out, nil
}

// replyMarkupMapping is a [pipeline.Mapping] that redirects a field typed as
// the four-variant reply markup union to ReplyMarkup; every other field
// rides through unchanged.
type replyMarkupMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (replyMarkupMapping) Apply(field typed.Field) (typed.Field, error) {
	if !field.Type.Equals(typeexpr.NewUnion(
		typeexpr.NewNamed(inlineKeyboardMarkupRef),
		typeexpr.NewNamed(replyKeyboardMarkupRef),
		typeexpr.NewNamed(replyKeyboardRemoveRef),
		typeexpr.NewNamed(forceReplyRef),
	)) {
		return field, nil
	}
	return typed.Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typeexpr.NewNamed(replyMarkupRef),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}
