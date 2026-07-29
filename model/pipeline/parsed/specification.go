// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
)

// Definitions represents the table of everything the page names — objects,
// methods, and unions alike — keyed by reference. A reference names at most
// one definition whatever its kind. Positions index the page's headings, so
// definitions of different kinds compare in the order the page presents them.
type Definitions = pipeline.Table[model.Reference, Definition]

// Fields represents the table of object fields, keyed by owner reference and
// field key. Positions are zero-based, unique and gapless within one owner's
// fields.
type Fields = pipeline.Table[model.FieldKey, Field]

// Params represents the table of method parameters, keyed by owner reference
// and parameter key. Positions are zero-based, unique and gapless within one
// owner's parameters.
type Params = pipeline.Table[model.FieldKey, Param]

// Variants represents the table of union variants, keyed by owner reference and
// variant reference. Positions are zero-based, unique and gapless within one
// union's variants.
type Variants = pipeline.Table[model.VariantKey, Variant]

// Specification is the initial database: tables transcribed straight from the
// page, before any interpretation. Fields and parameters are keyed by their
// owner's reference and their own key.
type Specification struct {
	Definitions Definitions
	Fields      Fields
	Params      Params
	Variants    Variants
	Release     Release
}

// Page is a parsed documentation page, the source of a Specification.
type Page struct {
	doc *goquery.Document
}

// NewPage constructs a Page over a parsed documentation page.
func NewPage(doc *goquery.Document) Page {
	return Page{doc: doc}
}

// Specification returns the database decoded from the page: the definition,
// field, parameter, and variant tables, plus the latest release. It fails when
// any section is malformed or when two definitions share a reference.
func (p Page) Specification() (Specification, error) {
	definitions, err := p.definitions()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding definitions: %w", err)
	}
	fields, err := NewFieldRows(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding fields: %w", err)
	}
	params, err := NewParamRows(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding parameters: %w", err)
	}
	variants, err := NewVariantItems(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding variants: %w", err)
	}
	release, err := NewChangelog(p.doc).Latest()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding release: %w", err)
	}
	return Specification{
		Definitions: definitions,
		Fields:      fields,
		Params:      params,
		Variants:    variants,
		Release:     release,
	}, nil
}

// definitions returns the object, method, and union sections gathered into one
// key space. It fails when any section is malformed or when two of them share a
// reference, since a reference names at most one definition.
func (p Page) definitions() (Definitions, error) {
	objects, err := NewObjectSections(p.doc).Table()
	if err != nil {
		return nil, fmt.Errorf("decoding objects: %w", err)
	}
	methods, err := NewMethodSections(p.doc).Table()
	if err != nil {
		return nil, fmt.Errorf("decoding methods: %w", err)
	}
	unions, err := NewUnionSections(p.doc).Table()
	if err != nil {
		return nil, fmt.Errorf("decoding unions: %w", err)
	}
	named, err := pipeline.NewMergedTable(objects, methods).Apply()
	if err != nil {
		return nil, fmt.Errorf("gathering object and method definitions: %w", err)
	}
	gathered, err := pipeline.NewMergedTable(named, unions).Apply()
	if err != nil {
		return nil, fmt.Errorf("gathering union definitions: %w", err)
	}
	return gathered, nil
}
