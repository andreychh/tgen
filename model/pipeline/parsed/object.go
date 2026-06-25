// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed/prose"
	prosetree "github.com/andreychh/tgen/model/prose"
)

// Object is the decoded record of a documentation object: its reference, name,
// and description. Its fields form a separate table.
type Object struct {
	Ref         model.Reference
	Name        model.Name
	Description prosetree.Tree
}

// ObjectSection is one object's section of the documentation page, headed by
// its <h4>.
type ObjectSection struct {
	h4 *goquery.Selection
}

// NewObjectSection constructs an ObjectSection over an object's <h4> header.
func NewObjectSection(h4 *goquery.Selection) ObjectSection {
	return ObjectSection{h4: h4}
}

// Record returns the object decoded from the section: its reference, name, and
// description. The description is the section's prose minus its field table, so
// notes after the table are kept. It fails when the reference, name, or
// description is malformed.
func (s ObjectSection) Record() (Object, error) {
	ref, err := NewReference(s.h4).Value()
	if err != nil {
		return Object{}, fmt.Errorf("parsing object reference: %w", err)
	}
	name, err := NewTypeName(s.h4).Value()
	if err != nil {
		return Object{}, fmt.Errorf("parsing object name: %w", err)
	}
	description, err := prose.NewParser(s.h4.NextUntil("h3, h4, hr").Not("table.table")).
		Parse()
	if err != nil {
		return Object{}, fmt.Errorf("parsing object description: %w", err)
	}
	return Object{
		Ref:         ref,
		Name:        name,
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

// Table returns the objects table, one record per object section. It fails when
// any object section is malformed.
func (s ObjectSections) Table() (pipeline.MapTable[model.Reference, Object], error) {
	out := pipeline.NewMapTable[model.Reference, Object]()
	for _, h4 := range s.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindObject {
			continue
		}
		object, err := NewObjectSection(h4).Record()
		if err != nil {
			return out, fmt.Errorf("parsing object: %w", err)
		}
		out.Insert(object.Ref, object)
	}
	return out, nil
}
