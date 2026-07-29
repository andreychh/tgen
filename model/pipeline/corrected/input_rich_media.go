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
	inputRichMediaRef      = model.Reference("inputrichmedia")
	inputMediaAnimationRef = model.Reference("inputmediaanimation")
	inputMediaAudioRef     = model.Reference("inputmediaaudio")
	inputMediaPhotoRef     = model.Reference("inputmediaphoto")
	inputMediaVideoRef     = model.Reference("inputmediavideo")
	inputMediaVoiceNoteRef = model.Reference("inputmediavoicenote")
)

// InputRichMedia is a [Rule] that introduces the InputRichMedia union — a
// media element embedded in a rich message — and redirects every field
// typed as their union to it.
type InputRichMedia struct{}

// Apply implements [Rule]. It fails when inputrichmedia already names
// something else in spec.
func (r InputRichMedia) Apply(spec Specification) (Specification, error) {
	definitions, err := r.definitions(spec.Definitions)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input rich media definitions: %w", err)
	}
	variants, err := r.variants(spec.Variants)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input rich media variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, inputRichMediaMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting input rich media fields: %w", err)
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

// definitions returns base holding what InputRichMedia names too: the union
// itself. Its variants are documented objects and name themselves. It fails
// when the reference is taken.
func (r InputRichMedia) definitions(base parsed.Definitions) (parsed.Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		inputRichMediaRef,
		"InputRichMedia",
		model.DefinitionKindUnion,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"InputRichMedia represents a media element embedded in a rich message.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the input rich media union: %w", err)
	}
	return out, nil
}

// variants returns base listing InputRichMedia's five variants, each already
// a documented object. It fails when the union already lists one of them.
func (r InputRichMedia) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, inputRichMediaRef)
	err := out.Insert(
		inputMediaAnimationRef,
		inputMediaAudioRef,
		inputMediaPhotoRef,
		inputMediaVideoRef,
		inputMediaVoiceNoteRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listing input rich media variants: %w", err)
	}
	return out, nil
}

// inputRichMediaMapping is a [pipeline.Mapping] that redirects a field
// typed as the five-variant input rich media union to InputRichMedia; every
// other field rides through unchanged.
type inputRichMediaMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (inputRichMediaMapping) Apply(field typed.Field) (typed.Field, error) {
	if !field.Type.Equals(typeexpr.NewUnion(
		typeexpr.NewNamed(inputMediaAnimationRef),
		typeexpr.NewNamed(inputMediaAudioRef),
		typeexpr.NewNamed(inputMediaPhotoRef),
		typeexpr.NewNamed(inputMediaVideoRef),
		typeexpr.NewNamed(inputMediaVoiceNoteRef),
	)) {
		return field, nil
	}
	return typed.Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typeexpr.NewNamed(inputRichMediaRef),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}
