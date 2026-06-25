// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
)

// Objects represents the table of standalone objects, keyed by reference.
type Objects = pipeline.Table[model.Reference, Object]

// Fields represents the table of object fields, keyed by owner reference and
// field key.
type Fields = pipeline.Table[model.FieldKey, Field]

// Methods represents the table of API methods, keyed by reference.
type Methods = pipeline.Table[model.Reference, Method]

// Params represents the table of method parameters, keyed by owner reference
// and parameter key.
type Params = pipeline.Table[model.FieldKey, Param]

// Unions represents the table of discriminated unions, keyed by reference.
type Unions = pipeline.Table[model.Reference, Union]

// Variants represents the table of union variants, keyed by owner reference and
// variant reference.
type Variants = pipeline.Table[model.VariantKey, Variant]

// Specification is the initial database: tables transcribed straight from the
// page, before any interpretation. Fields and parameters are keyed by their
// owner's reference and their own key.
type Specification struct {
	Objects  Objects
	Fields   Fields
	Methods  Methods
	Params   Params
	Unions   Unions
	Variants Variants
	Release  Release
}

// Page is a parsed documentation page, the source of a Specification.
type Page struct {
	doc *goquery.Document
}

// NewPage constructs a Page over a parsed documentation page.
func NewPage(doc *goquery.Document) Page {
	return Page{doc: doc}
}

// Specification returns the database decoded from the page: the object, field,
// method, parameter, union, and variant tables, plus the latest release. It
// fails when any section is malformed.
func (p Page) Specification() (Specification, error) {
	objects, err := NewObjectSections(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding objects: %w", err)
	}
	fields, err := NewFieldRows(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding fields: %w", err)
	}
	methods, err := NewMethodSections(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding methods: %w", err)
	}
	params, err := NewParamRows(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding parameters: %w", err)
	}
	unions, err := NewUnionSections(p.doc).Table()
	if err != nil {
		return Specification{}, fmt.Errorf("decoding unions: %w", err)
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
		Objects:  objects,
		Fields:   fields,
		Methods:  methods,
		Params:   params,
		Unions:   unions,
		Variants: variants,
		Release:  release,
	}, nil
}
