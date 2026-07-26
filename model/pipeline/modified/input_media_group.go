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
	if alreadyExists(spec, inputMediaGroupRef) {
		return Specification{}, fmt.Errorf(
			"%s already names something else in spec",
			inputMediaGroupRef,
		)
	}
	unions, err := pipeline.NewMergedTable(spec.Unions, r.union()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input media group union: %w", err)
	}
	variants, err := pipeline.NewMergedTable(spec.Variants, r.variants()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input media group variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, inputMediaGroupMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting input media group fields: %w", err)
	}
	return Specification{
		Objects:  spec.Objects,
		Methods:  spec.Methods,
		Fields:   fields,
		Unions:   unions,
		Variants: variants,
		Aliases:  spec.Aliases,
		Release:  spec.Release,
	}, nil
}

// union returns the table holding the single InputMediaGroup union
// definition.
func (r InputMediaGroup) union() parsed.Unions {
	out := pipeline.NewMapTableWithCapacity[model.Reference, parsed.Union](1)
	out.Insert(inputMediaGroupRef, parsed.Union{
		Ref:  inputMediaGroupRef,
		Name: "InputMediaGroup",
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"InputMediaGroup represents a media element in a media group.",
			prose.StylePlain,
		))),
	})
	return out
}

// variants returns the table of InputMediaGroup's five variants, each
// already a documented object.
func (r InputMediaGroup) variants() parsed.Variants {
	out := pipeline.NewMapTableWithCapacity[model.VariantKey, parsed.Variant](5)
	out.Insert(
		model.VariantKey{Owner: inputMediaGroupRef, Ref: inputMediaAudioRef},
		parsed.Variant{Ref: inputMediaAudioRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputMediaGroupRef, Ref: inputMediaDocumentRef},
		parsed.Variant{Ref: inputMediaDocumentRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputMediaGroupRef, Ref: inputMediaLivePhotoRef},
		parsed.Variant{Ref: inputMediaLivePhotoRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputMediaGroupRef, Ref: inputMediaPhotoRef},
		parsed.Variant{Ref: inputMediaPhotoRef},
	)
	out.Insert(
		model.VariantKey{Owner: inputMediaGroupRef, Ref: inputMediaVideoRef},
		parsed.Variant{Ref: inputMediaVideoRef},
	)
	return out
}

// inputMediaGroupMapping is a [pipeline.Mapping] that redirects a field
// typed as an array of the five-variant input media union to an array of
// InputMediaGroup; every other field rides through unchanged.
type inputMediaGroupMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (inputMediaGroupMapping) Apply(field typed.Field) (typed.Field, error) {
	if !field.Type.Equals(typetree.NewArray(typetree.NewUnion(
		typetree.NewNamed(inputMediaAudioRef),
		typetree.NewNamed(inputMediaDocumentRef),
		typetree.NewNamed(inputMediaLivePhotoRef),
		typetree.NewNamed(inputMediaPhotoRef),
		typetree.NewNamed(inputMediaVideoRef),
	))) {
		return field, nil
	}
	return typed.Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typetree.NewArray(typetree.NewNamed(inputMediaGroupRef)),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}
