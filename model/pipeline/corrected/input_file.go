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

const (
	inputFileRef = model.Reference("inputfile")
	fileIDRef    = model.Reference("fileid")
	uploadRef    = model.Reference("upload")
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
// taking the place of the documented object base no longer holds, and the file
// id alias its variant stands for. It fails when either reference is taken.
func (r InputFile) definitions(base parsed.Definitions) (parsed.Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		inputFileRef,
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
	return out, nil
}

// aliases returns the type the file id alias stands for: the documentation
// writes it as a plain String field, not a named type.
func (r InputFile) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](1)
	out.Insert(fileIDRef, Alias{Ref: fileIDRef, Type: typeexpr.NewPrimitive(primitive.String)})
	return out
}

// variants returns base listing InputFile's variants, which is FileId alone:
// the upload variant holds a Go io.Reader with no shape shared across targets,
// so it still needs its own design before it can join the union. It fails when
// the union already lists FileId.
func (r InputFile) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, inputFileRef)
	err := out.Insert(
		fileIDRef,
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
func (inputFileFilter) Apply(ref model.Reference, definition parsed.Definition) bool {
	return ref != inputFileRef
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
		Type:        typeexpr.NewNamed(inputFileRef),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}

// matches reports whether field is typed as InputFile alone, as InputFile or
// String, or as a bare String field whose description links to the
// sending-files guide.
func (m inputFileMapping) matches(field typed.Field) bool {
	return field.Type.Equals(typeexpr.NewNamed(inputFileRef)) ||
		field.Type.Equals(typeexpr.NewUnion(
			typeexpr.NewNamed(inputFileRef),
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
