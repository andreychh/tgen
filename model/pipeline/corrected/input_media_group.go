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
	inputMediaGroupRef     = model.Reference("inputmediagroup")
	inputMediaDocumentRef  = model.Reference("inputmediadocument")
	inputMediaLivePhotoRef = model.Reference("inputmedialivephoto")
)

// InputMediaGroup is a [Rule] that introduces the InputMediaGroup union — a
// media element in a media group — and redirects every field typed as an
// array of their union to an array of it.
type InputMediaGroup struct{}

// Apply implements [Rule]. It fails when inputmediagroup already names
// something else in spec.
func (r InputMediaGroup) Apply(spec Specification) (Specification, error) {
	definitions, err := r.definitions(spec.Definitions)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input media group definitions: %w", err)
	}
	variants, err := r.variants(spec.Variants)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input media group variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, inputMediaGroupMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting input media group fields: %w", err)
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

// definitions returns base holding what InputMediaGroup names too: the union
// itself. Its variants are documented objects and name themselves. It fails
// when the reference is taken.
func (r InputMediaGroup) definitions(base Definitions) (Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		inputMediaGroupRef,
		"InputMediaGroup",
		model.DefinitionKindUnion,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"InputMediaGroup represents a media element in a media group.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the input media group union: %w", err)
	}
	return out, nil
}

// variants returns base listing InputMediaGroup's five variants, each already
// a documented object. It fails when the union already lists one of them.
func (r InputMediaGroup) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, inputMediaGroupRef)
	err := out.Insert(
		inputMediaAudioRef,
		inputMediaDocumentRef,
		inputMediaLivePhotoRef,
		inputMediaPhotoRef,
		inputMediaVideoRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listing input media group variants: %w", err)
	}
	return out, nil
}

// inputMediaGroupMapping is a [pipeline.Mapping] that redirects a field
// typed as an array of the five-variant input media union to an array of
// InputMediaGroup; every other field rides through unchanged.
type inputMediaGroupMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (inputMediaGroupMapping) Apply(field typed.Field) (typed.Field, error) {
	if !field.Type.Equals(typeexpr.NewArray(typeexpr.NewUnion(
		typeexpr.NewNamed(inputMediaAudioRef),
		typeexpr.NewNamed(inputMediaDocumentRef),
		typeexpr.NewNamed(inputMediaLivePhotoRef),
		typeexpr.NewNamed(inputMediaPhotoRef),
		typeexpr.NewNamed(inputMediaVideoRef),
	))) {
		return field, nil
	}
	return typed.Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typeexpr.NewArray(typeexpr.NewNamed(inputMediaGroupRef)),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}
