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

// Method is the decoded record of a documentation method: its reference, name,
// and description. Its parameters form a separate table.
type Method struct {
	Ref         model.Reference
	Name        model.Name
	Description prosetree.Tree
}

// MethodSection is one method's section of the documentation page, headed by
// its <h4>.
type MethodSection struct {
	h4 *goquery.Selection
}

// NewMethodSection constructs a MethodSection over a method's <h4> header.
func NewMethodSection(h4 *goquery.Selection) MethodSection {
	return MethodSection{h4: h4}
}

// Record returns the method decoded from the section: its reference, name, and
// description. The return type stays inside the description prose for a later
// pass to interpret. It fails when the reference, name, or description is
// malformed.
func (s MethodSection) Record() (Method, error) {
	ref, err := NewReference(s.h4).Value()
	if err != nil {
		return Method{}, fmt.Errorf("parsing method reference: %w", err)
	}
	name, err := NewMethodName(s.h4).Value()
	if err != nil {
		return Method{}, fmt.Errorf("parsing method name: %w", err)
	}
	description, err := prose.NewParser(s.h4.NextUntil("h3, h4, hr").Not("table.table")).
		Parse()
	if err != nil {
		return Method{}, fmt.Errorf("parsing method description: %w", err)
	}
	return Method{
		Ref:         ref,
		Name:        name,
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

// Table returns the methods table, one record per method section. It fails when
// any method section is malformed.
func (s MethodSections) Table() (pipeline.MapTable[model.Reference, Method], error) {
	out := pipeline.NewMapTable[model.Reference, Method]()
	for _, h4 := range s.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindMethod {
			continue
		}
		method, err := NewMethodSection(h4).Record()
		if err != nil {
			return out, fmt.Errorf("parsing method: %w", err)
		}
		out.Insert(method.Ref, method)
	}
	return out, nil
}
