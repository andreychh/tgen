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

// ObjectSection is one object's section of the documentation page, headed by
// its <h4> at its position among the page's headings.
type ObjectSection struct {
	at int
	h4 *goquery.Selection
}

// NewObjectSection constructs an ObjectSection over an object's <h4> header at
// position at.
func NewObjectSection(at int, h4 *goquery.Selection) ObjectSection {
	return ObjectSection{at: at, h4: h4}
}

// Record returns the object decoded from the section as a definition of object
// kind: its reference, name, position, and description. The description is the
// section's prose minus its field table, so notes after the table are kept. It
// fails when the reference, name, or description is malformed.
func (s ObjectSection) Record() (Definition, error) {
	ref, err := NewReference(s.h4).Value()
	if err != nil {
		return Definition{}, fmt.Errorf("parsing object reference: %w", err)
	}
	name, err := NewTypeName(s.h4).Value()
	if err != nil {
		return Definition{}, fmt.Errorf("parsing object name: %w", err)
	}
	description, err := prose.NewPassage(s.h4.NextUntil("h3, h4, hr").Not("table.table")).Value()
	if err != nil {
		return Definition{}, fmt.Errorf("parsing object description: %w", err)
	}
	return Definition{
		Ref:         ref,
		Name:        name,
		Kind:        model.DefinitionKindObject,
		Position:    model.Position(s.at),
		Description: description,
	}, nil
}

// ObjectSections are the object sections of a documentation page.
type ObjectSections struct {
	doc *goquery.Document
}

// NewObjectSections constructs an ObjectSections over a parsed documentation
// page.
func NewObjectSections(doc *goquery.Document) ObjectSections {
	return ObjectSections{doc: doc}
}

// Table returns the definitions of object kind, one record per object section,
// each holding the position of its heading among the page's headings. It fails
// when any object section is malformed.
func (s ObjectSections) Table() (pipeline.MapTable[model.Reference, Definition], error) {
	out := pipeline.NewMapTable[model.Reference, Definition]()
	for at, h4 := range s.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindObject {
			continue
		}
		object, err := NewObjectSection(at, h4).Record()
		if err != nil {
			return out, fmt.Errorf("parsing object: %w", err)
		}
		out.Insert(object.Ref, object)
	}
	return out, nil
}
