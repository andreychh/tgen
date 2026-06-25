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

// Field is the decoded record of one object field: its key and the verbatim
// prose of its type and description. Resolving the type and lifting optionality
// are left for later passes.
type Field struct {
	Key         model.Key
	Type        prosetree.Phrase
	Description prosetree.Phrase
}

// FieldRow is one field's row of an object's table.
type FieldRow struct {
	tr *goquery.Selection
}

// NewFieldRow constructs a FieldRow over a table row.
func NewFieldRow(tr *goquery.Selection) FieldRow {
	return FieldRow{tr: tr}
}

// Record returns the field decoded from the row: its key and the prose of its
// type and description. It fails when the key, type, or description is
// malformed.
func (r FieldRow) Record() (Field, error) {
	key, err := NewKey(r.cell(0)).Value()
	if err != nil {
		return Field{}, fmt.Errorf("parsing field key: %w", err)
	}
	typ, err := prose.NewPhrase(r.cell(1).Contents()).Value()
	if err != nil {
		return Field{}, fmt.Errorf("parsing field type: %w", err)
	}
	description, err := prose.NewPhrase(r.cell(2).Contents()).Value()
	if err != nil {
		return Field{}, fmt.Errorf("parsing field description: %w", err)
	}
	return Field{
		Key:         key,
		Type:        typ,
		Description: description,
	}, nil
}

func (r FieldRow) cell(index int) *goquery.Selection {
	return r.tr.ChildrenFiltered("td").Eq(index)
}

// ObjectFields are the field rows declared under one object's heading.
type ObjectFields struct {
	h4 *goquery.Selection
}

// NewObjectFields constructs an ObjectFields over an object's <h4> header.
func NewObjectFields(h4 *goquery.Selection) ObjectFields {
	return ObjectFields{h4: h4}
}

// Records returns the fields under the heading, paired with the owning object
// reference. It fails when the reference or any field row is malformed.
func (f ObjectFields) Records() (model.Reference, []Field, error) {
	owner, err := NewReference(f.h4).Value()
	if err != nil {
		return "", nil, fmt.Errorf("parsing object reference: %w", err)
	}
	var fields []Field
	rows := f.h4.NextUntil("h3, h4, hr").Filter("table.table").First().Find("tbody > tr")
	for _, tr := range rows.EachIter() {
		field, err := NewFieldRow(tr).Record()
		if err != nil {
			return "", nil, fmt.Errorf("parsing field: %w", err)
		}
		fields = append(fields, field)
	}
	return owner, fields, nil
}

// FieldRows are the field rows of every object on a documentation page.
type FieldRows struct {
	doc *goquery.Document
}

// NewFieldRows constructs a FieldRows over a parsed documentation page.
func NewFieldRows(doc *goquery.Document) FieldRows {
	return FieldRows{doc: doc}
}

// Table returns the fields table, one record per field row, keyed by owning
// object and field key. It fails when any reference or field row is malformed.
func (r FieldRows) Table() (pipeline.MapTable[model.FieldKey, Field], error) {
	out := pipeline.NewMapTable[model.FieldKey, Field]()
	for _, h4 := range r.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindObject {
			continue
		}
		owner, fields, err := NewObjectFields(h4).Records()
		if err != nil {
			return out, err
		}
		for _, field := range fields {
			out.Insert(model.FieldKey{Owner: owner, Key: field.Key}, field)
		}
	}
	return out, nil
}
