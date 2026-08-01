// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/typed"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

// InputFileRef is the reference of the InputFile union this stage introduces.
// A later stage decides which definitions can carry a file, and it starts from
// here: InputFile is the only type that is one, and the documentation never
// names it, so nothing downstream can find it on its own.
const InputFileRef = model.Reference("inputfile")

const (
	fileIDRef = model.Reference("fileid")
	uploadRef = model.Reference("upload")
)

// InputFile is a [Rule] that introduces the InputFile union — a file to
// send, either by file ID or by uploading — in place of the documented
// InputFile object, and redirects every field that could carry one to it.
type InputFile struct{}

// Apply implements [Rule]. It fails when fileid already names something else
// in spec.
func (r InputFile) Apply(spec Specification) (Specification, error) {
	documented := pipeline.NewFilteredTable(spec.Definitions, inputFileFilter{}).Apply()
	definitions, err := r.definitions(documented)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input file definitions: %w", err)
	}
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input file aliases: %w", err)
	}
	variants, err := r.variants(spec.Variants)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input file variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, inputFileMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting input file fields: %w", err)
	}
	return Specification{
		Definitions:    definitions,
		Methods:        spec.Methods,
		Fields:         fields,
		Discriminators: spec.Discriminators,
		Variants:       variants,
		Aliases:        aliases,
		Release:        spec.Release,
	}, nil
}

// definitions returns base holding what InputFile names too: the union itself,
// taking the place of the documented object base no longer holds, the file id
// alias one variant stands for, and the upload object the other is. The upload
// object owns no field — what it holds, the bytes to send and the name to send
// them under, has no shape shared across targets, so each target spells it
// itself. It fails when any of the three references is taken.
func (r InputFile) definitions(base Definitions) (Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		InputFileRef,
		"InputFile",
		model.DefinitionKindUnion,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"InputFile represents a file to send, either by file ID or by uploading.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the input file union: %w", err)
	}
	err = out.Insert(
		fileIDRef,
		"FileId",
		model.DefinitionKindAlias,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"FileID represents a Telegram file identifier.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the file id alias: %w", err)
	}
	err = out.Insert(
		uploadRef,
		"Upload",
		model.DefinitionKindObject,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"Upload represents a file sent with the request, carrying the bytes to send "+
				"and the name to send them under.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the upload object: %w", err)
	}
	return out, nil
}

// aliases returns the type the file id alias stands for: the documentation
// writes it as a plain String field, not a named type.
func (r InputFile) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](1)
	out.Insert(fileIDRef, Alias{Ref: fileIDRef, Type: typeexpr.NewPrimitive(primitive.String)})
	return out
}

// variants returns base listing InputFile's variants: a file Telegram already
// holds, named by its id, and one uploaded with the request. It fails when the
// union already lists either.
func (r InputFile) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, InputFileRef)
	err := out.Insert(
		fileIDRef,
		uploadRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listing input file variants: %w", err)
	}
	return out, nil
}

// inputFileFilter is a [pipeline.Filter] that excludes the documented
// InputFile object, whose place the InputFile union takes.
type inputFileFilter struct{}

// Apply implements [pipeline.Filter]. It reports whether definition belongs in
// the filtered result: false when ref is inputfile.
func (inputFileFilter) Apply(ref model.Reference, definition Definition) bool {
	return ref != InputFileRef
}

// inputFileMapping is a [pipeline.Mapping] that redirects a field that could
// carry an InputFile to InputFile; every other field rides through
// unchanged.
type inputFileMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (m inputFileMapping) Apply(field typed.Field) (typed.Field, error) {
	if !m.matches(field) {
		return field, nil
	}
	return typed.Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typeexpr.NewNamed(InputFileRef),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}

// matches reports whether field is typed as InputFile alone, as InputFile or
// String, or as a bare String field whose description links to the
// sending-files guide.
func (m inputFileMapping) matches(field typed.Field) bool {
	return field.Type.Equals(typeexpr.NewNamed(InputFileRef)) ||
		field.Type.Equals(typeexpr.NewUnion(
			typeexpr.NewNamed(InputFileRef),
			typeexpr.NewPrimitive(primitive.String),
		)) ||
		field.Type.Equals(typeexpr.NewPrimitive(primitive.String)) &&
			hasSendingFilesLink(field.Description)
}

// hasSendingFilesLink reports whether description links to the sending-files
// guide section.
func hasSendingFilesLink(description prose.Phrase) bool {
	for _, inline := range description.Inlines() {
		link, ok := inline.(prose.Link)
		if !ok {
			continue
		}
		anchor, isAnchor := link.Anchor()
		if isAnchor && anchor == "sending-files" {
			return true
		}
	}
	return false
}
