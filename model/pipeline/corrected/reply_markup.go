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
	if alreadyExists(spec, replyMarkupRef) {
		return Specification{}, fmt.Errorf(
			"%s already names something else in spec",
			replyMarkupRef,
		)
	}
	unions, err := pipeline.NewMergedTable(spec.Unions, r.union()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing reply markup union: %w", err)
	}
	variants, err := pipeline.NewMergedTable(spec.Variants, r.variants()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing reply markup variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, replyMarkupMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting reply markup fields: %w", err)
	}
	return Specification{
		Objects:        spec.Objects,
		Methods:        spec.Methods,
		Fields:         fields,
		Discriminators: spec.Discriminators,
		Unions:         unions,
		Variants:       variants,
		Aliases:        spec.Aliases,
		Release:        spec.Release,
	}, nil
}

// union returns the table holding the single ReplyMarkup union definition.
func (r ReplyMarkup) union() parsed.Unions {
	out := pipeline.NewMapTableWithCapacity[model.Reference, parsed.Union](1)
	out.Insert(replyMarkupRef, parsed.Union{
		Ref:  replyMarkupRef,
		Name: "ReplyMarkup",
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"ReplyMarkup represents a reply markup attached to a message.",
			prose.StylePlain,
		))),
	})
	return out
}

// variants returns the table of ReplyMarkup's four variants, each already a
// documented object.
func (r ReplyMarkup) variants() parsed.Variants {
	out := pipeline.NewMapTableWithCapacity[model.VariantKey, parsed.Variant](4)
	out.Insert(
		model.VariantKey{Owner: replyMarkupRef, Ref: inlineKeyboardMarkupRef},
		parsed.Variant{Ref: inlineKeyboardMarkupRef},
	)
	out.Insert(
		model.VariantKey{Owner: replyMarkupRef, Ref: replyKeyboardMarkupRef},
		parsed.Variant{Ref: replyKeyboardMarkupRef},
	)
	out.Insert(
		model.VariantKey{Owner: replyMarkupRef, Ref: replyKeyboardRemoveRef},
		parsed.Variant{Ref: replyKeyboardRemoveRef},
	)
	out.Insert(
		model.VariantKey{Owner: replyMarkupRef, Ref: forceReplyRef},
		parsed.Variant{Ref: forceReplyRef},
	)
	return out
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
