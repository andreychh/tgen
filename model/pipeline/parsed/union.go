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

// Union is the decoded record of a documentation union: its reference, name,
// and description. Its variants form a separate table.
type Union struct {
	Ref         model.Reference
	Name        model.Name
	Description prosetree.Passage
}

// UnionSection is one union's section of the documentation page, headed by its
// <h4>.
type UnionSection struct {
	h4 *goquery.Selection
}

// NewUnionSection constructs a UnionSection over a union's <h4> header.
func NewUnionSection(h4 *goquery.Selection) UnionSection {
	return UnionSection{h4: h4}
}

// Record returns the union decoded from the section: its reference, name, and
// description. The description is the section's prose minus its variant list.
// It fails when the reference, name, or description is malformed.
func (s UnionSection) Record() (Union, error) {
	ref, err := NewReference(s.h4).Value()
	if err != nil {
		return Union{}, fmt.Errorf("parsing union reference: %w", err)
	}
	name, err := NewTypeName(s.h4).Value()
	if err != nil {
		return Union{}, fmt.Errorf("parsing union name: %w", err)
	}
	description, err := prose.NewPassage(s.h4.NextUntil("h3, h4, hr").Not("ul")).Value()
	if err != nil {
		return Union{}, fmt.Errorf("parsing union description: %w", err)
	}
	return Union{
		Ref:         ref,
		Name:        name,
		Description: description,
	}, nil
}

// UnionSections are the union sections of a documentation page.
type UnionSections struct {
	doc *goquery.Document
}

// NewUnionSections constructs a UnionSections over a parsed documentation page.
func NewUnionSections(doc *goquery.Document) UnionSections {
	return UnionSections{doc: doc}
}

// Table returns the unions table, one record per union section. It fails when
// any union section is malformed.
func (s UnionSections) Table() (pipeline.MapTable[model.Reference, Union], error) {
	out := pipeline.NewMapTable[model.Reference, Union]()
	for _, h4 := range s.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindUnion {
			continue
		}
		union, err := NewUnionSection(h4).Record()
		if err != nil {
			return out, fmt.Errorf("parsing union: %w", err)
		}
		out.Insert(union.Ref, union)
	}
	return out, nil
}
