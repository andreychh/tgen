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
	inputFileRef = model.Reference("inputfile")
	fileIDRef    = model.Reference("fileid")
	uploadRef    = model.Reference("upload")
)

// InputFile is a [Rule] that introduces the InputFile union — a file to
// send, either by file ID or by uploading — in place of the documented
// InputFile object, and redirects every field that could carry one to it.
type InputFile struct{}

// Apply implements [Rule]. It fails when fileid or upload already names
// something else in spec.
func (r InputFile) Apply(spec Specification) (Specification, error) {
	for _, ref := range []model.Reference{fileIDRef, uploadRef} {
		if alreadyExists(spec, ref) {
			return Specification{}, fmt.Errorf("%s already names something else in spec", ref)
		}
	}
	objects := pipeline.NewFilteredTable(spec.Objects, inputFileFilter{}).Apply()
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input file aliases: %w", err)
	}
	unions, err := pipeline.NewMergedTable(spec.Unions, r.union()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input file union: %w", err)
	}
	variants, err := pipeline.NewMergedTable(spec.Variants, r.variants()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing input file variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, inputFileMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting input file fields: %w", err)
	}
	return Specification{
		Objects:  objects,
		Methods:  spec.Methods,
		Fields:   fields,
		Unions:   unions,
		Variants: variants,
		Aliases:  aliases,
		Release:  spec.Release,
	}, nil
}

// aliases returns the table of primitive aliases InputFile's variants need:
// a file id, which the documentation writes as a plain String field, not a
// named type.
func (r InputFile) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](1)
	out.Insert(fileIDRef, Alias{
		Ref:  fileIDRef,
		Name: "FileId",
		Type: typetree.NewPrimitive(typetree.String),
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"FileID represents a Telegram file identifier.",
			prose.StylePlain,
		))),
	})
	return out
}

// union returns the table holding the single InputFile union definition.
func (r InputFile) union() parsed.Unions {
	out := pipeline.NewMapTableWithCapacity[model.Reference, parsed.Union](1)
	out.Insert(inputFileRef, parsed.Union{
		Ref:  inputFileRef,
		Name: "InputFile",
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"InputFile represents a file to send, either by file ID or by uploading.",
			prose.StylePlain,
		))),
	})
	return out
}

// variants returns the table of InputFile's variants. It carries only
// FileId: the upload variant holds a Go io.Reader with no shape shared
// across targets, so it still needs its own design before it can join this
// table.
func (r InputFile) variants() parsed.Variants {
	out := pipeline.NewMapTableWithCapacity[model.VariantKey, parsed.Variant](1)
	out.Insert(
		model.VariantKey{Owner: inputFileRef, Ref: fileIDRef},
		parsed.Variant{Ref: fileIDRef},
	)
	return out
}

// inputFileFilter is a [pipeline.Filter] that excludes the documented InputFile
// object from Objects.
type inputFileFilter struct{}

// Apply implements [pipeline.Filter]. It reports whether object belongs in the
// filtered result: false when ref is inputfile.
func (inputFileFilter) Apply(ref model.Reference, object parsed.Object) bool {
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
		Type:        typetree.NewNamed(inputFileRef),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}

// matches reports whether field is typed as InputFile alone, as InputFile or
// String, or as a bare String field whose description links to the
// sending-files guide.
func (m inputFileMapping) matches(field typed.Field) bool {
	return field.Type.Equals(typetree.NewNamed(inputFileRef)) ||
		field.Type.Equals(typetree.NewUnion(
			typetree.NewNamed(inputFileRef),
			typetree.NewPrimitive(typetree.String),
		)) ||
		field.Type.Equals(typetree.NewPrimitive(typetree.String)) &&
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
