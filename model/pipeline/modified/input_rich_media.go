// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package modified

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/typed"
	"github.com/andreychh/tgen/model/prose"
	typetree "github.com/andreychh/tgen/model/types/v2"
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
	if alreadyExists(spec, inputRichMediaRef) {
		return Specification{}, fmt.Errorf(
			"%s already names something else in spec",
			inputRichMediaRef,
		)
	}
	unions, err := pipeline.NewMergedTable(spec.Unions, r.union()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input rich media union: %w", err)
	}
	variants, err := pipeline.NewMergedTable(spec.Variants, r.variants()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input rich media variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, inputRichMediaMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting input rich media fields: %w", err)
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

// union returns the table holding the single InputRichMedia union
// definition.
func (r InputRichMedia) union() parsed.Unions {
	out := pipeline.NewMapTableWithCapacity[model.Reference, parsed.Union](1)
	out.Insert(inputRichMediaRef, parsed.Union{
		Ref:  inputRichMediaRef,
		Name: "InputRichMedia",
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"InputRichMedia represents a media element embedded in a rich message.",
			prose.StylePlain,
		))),
	})
	return out
}

// variants returns the table of InputRichMedia's five variants, each
// already a documented object.
func (r InputRichMedia) variants() parsed.Variants {
	out := pipeline.NewMapTableWithCapacity[model.VariantKey, parsed.Variant](5)
	out.Insert(
		model.VariantKey{Owner: inputRichMediaRef, Ref: inputMediaAnimationRef},
		parsed.Variant{Ref: inputMediaAnimationRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputRichMediaRef, Ref: inputMediaAudioRef},
		parsed.Variant{Ref: inputMediaAudioRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputRichMediaRef, Ref: inputMediaPhotoRef},
		parsed.Variant{Ref: inputMediaPhotoRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputRichMediaRef, Ref: inputMediaVideoRef},
		parsed.Variant{Ref: inputMediaVideoRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputRichMediaRef, Ref: inputMediaVoiceNoteRef},
		parsed.Variant{Ref: inputMediaVoiceNoteRef},
	)
	return out
}

// inputRichMediaMapping is a [pipeline.Mapping] that redirects a field
// typed as the five-variant input rich media union to InputRichMedia; every
// other field rides through unchanged.
type inputRichMediaMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (inputRichMediaMapping) Apply(field typed.Field) (typed.Field, error) {
	if !field.Type.Equals(typetree.NewUnion(
		typetree.NewNamed(inputMediaAnimationRef),
		typetree.NewNamed(inputMediaAudioRef),
		typetree.NewNamed(inputMediaPhotoRef),
		typetree.NewNamed(inputMediaVideoRef),
		typetree.NewNamed(inputMediaVoiceNoteRef),
	)) {
		return field, nil
	}
	return typed.Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typetree.NewNamed(inputRichMediaRef),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}
