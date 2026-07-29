// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed/prose"
)

// MethodSection is one method's section of the documentation page, headed by
// its <h4> at its position among the page's headings.
type MethodSection struct {
	at int
	h4 *goquery.Selection
}

// NewMethodSection constructs a MethodSection over a method's <h4> header at
// position at.
func NewMethodSection(at int, h4 *goquery.Selection) MethodSection {
	return MethodSection{at: at, h4: h4}
}

// Record returns the method decoded from the section as a definition of method
// kind: its reference, name, position, and description. The return type stays
// inside the description prose for a later pass to interpret. It fails when the
// reference, name, or description is malformed.
func (s MethodSection) Record() (Definition, error) {
	ref, err := NewReference(s.h4).Value()
	if err != nil {
		return Definition{}, fmt.Errorf("parsing method reference: %w", err)
	}
	name, err := NewMethodName(s.h4).Value()
	if err != nil {
		return Definition{}, fmt.Errorf("parsing method name: %w", err)
	}
	description, err := prose.NewPassage(s.h4.NextUntil("h3, h4, hr").Not("table.table")).Value()
	if err != nil {
		return Definition{}, fmt.Errorf("parsing method description: %w", err)
	}
	return Definition{
		Ref:         ref,
		Name:        name,
		Kind:        model.DefinitionKindMethod,
		Position:    model.Position(s.at),
		Description: description,
	}, nil
}

// MethodSections are the method sections of a documentation page.
type MethodSections struct {
	doc *goquery.Document
}

// NewMethodSections constructs a MethodSections over a parsed documentation
// page.
func NewMethodSections(doc *goquery.Document) MethodSections {
	return MethodSections{doc: doc}
}

// Table returns the definitions of method kind, one record per method section,
// each holding the position of its heading among the page's headings. It fails
// when any method section is malformed.
func (s MethodSections) Table() (pipeline.MapTable[model.Reference, Definition], error) {
	out := pipeline.NewMapTable[model.Reference, Definition]()
	for at, h4 := range s.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindMethod {
			continue
		}
		method, err := NewMethodSection(at, h4).Record()
		if err != nil {
			return out, fmt.Errorf("parsing method: %w", err)
		}
		out.Insert(method.Ref, method)
	}
	return out, nil
}
