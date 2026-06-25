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
	Type        prosetree.Tree
	Description prosetree.Tree
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
	typ, err := prose.NewParser(r.cell(1).Contents()).Parse()
	if err != nil {
		return Field{}, fmt.Errorf("parsing field type: %w", err)
	}
	description, err := prose.NewParser(r.cell(2).Contents()).Parse()
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
		ref, err := NewReference(h4).Value()
		if err != nil {
			return out, fmt.Errorf("parsing object reference: %w", err)
		}
		rows := h4.NextUntil("h3, h4, hr").Filter("table.table").First().Find("tbody > tr")
		for _, tr := range rows.EachIter() {
			field, err := NewFieldRow(tr).Record()
			if err != nil {
				return out, fmt.Errorf("parsing field: %w", err)
			}
			out.Insert(model.FieldKey{Owner: ref, Key: field.Key}, field)
		}
	}
	return out, nil
}
