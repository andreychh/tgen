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

// Param is the decoded record of one method parameter: its key, its position
// among the method's other parameters, and the verbatim prose of its type,
// requiredness, and description. Lifting optionality is left for a later
// pass, as with object fields.
type Param struct {
	Key         model.Key
	Position    model.Position
	Type        prosetree.Phrase
	Required    prosetree.Phrase
	Description prosetree.Phrase
}

// ParamRow is one parameter's row of a method's table, at its position among
// sibling rows.
type ParamRow struct {
	at int
	tr *goquery.Selection
}

// NewParamRow constructs a ParamRow over a table row at position at.
func NewParamRow(at int, tr *goquery.Selection) ParamRow {
	return ParamRow{at: at, tr: tr}
}

// Record returns the parameter decoded from the row: its key, its position,
// and the prose of its type, requiredness, and description. It fails when the
// key, type, requiredness, or description is malformed.
func (r ParamRow) Record() (Param, error) {
	key, err := NewKey(r.cell(0)).Value()
	if err != nil {
		return Param{}, fmt.Errorf("parsing parameter key: %w", err)
	}
	typ, err := prose.NewPhrase(r.cell(1).Contents()).Value()
	if err != nil {
		return Param{}, fmt.Errorf("parsing parameter type: %w", err)
	}
	required, err := prose.NewPhrase(r.cell(2).Contents()).Value()
	if err != nil {
		return Param{}, fmt.Errorf("parsing parameter required: %w", err)
	}
	description, err := prose.NewPhrase(r.cell(3).Contents()).Value()
	if err != nil {
		return Param{}, fmt.Errorf("parsing parameter description: %w", err)
	}
	return Param{
		Key:         key,
		Position:    model.Position(r.at),
		Type:        typ,
		Required:    required,
		Description: description,
	}, nil
}

func (r ParamRow) cell(index int) *goquery.Selection {
	return r.tr.ChildrenFiltered("td").Eq(index)
}

// MethodParams are the parameter rows declared under one method's heading.
type MethodParams struct {
	h4 *goquery.Selection
}

// NewMethodParams constructs a MethodParams over a method's <h4> header.
func NewMethodParams(h4 *goquery.Selection) MethodParams {
	return MethodParams{h4: h4}
}

// Records returns the parameters under the heading, paired with the owning
// method reference. It fails when the reference or any parameter row is
// malformed.
func (p MethodParams) Records() (model.Reference, []Param, error) {
	owner, err := NewReference(p.h4).Value()
	if err != nil {
		return "", nil, fmt.Errorf("parsing method reference: %w", err)
	}
	var params []Param
	rows := p.h4.NextUntil("h3, h4, hr").Filter("table.table").First().Find("tbody > tr")
	for at, tr := range rows.EachIter() {
		param, err := NewParamRow(at, tr).Record()
		if err != nil {
			return "", nil, fmt.Errorf("parsing parameter: %w", err)
		}
		params = append(params, param)
	}
	return owner, params, nil
}

// ParamRows are the parameter rows of every method on a documentation page.
type ParamRows struct {
	doc *goquery.Document
}

// NewParamRows constructs a ParamRows over a parsed documentation page.
func NewParamRows(doc *goquery.Document) ParamRows {
	return ParamRows{doc: doc}
}

// Table returns the parameters table, one record per parameter row, keyed by
// owning method and parameter key. It fails when any reference or parameter row
// is malformed.
func (r ParamRows) Table() (pipeline.MapTable[model.FieldKey, Param], error) {
	out := pipeline.NewMapTable[model.FieldKey, Param]()
	for _, h4 := range r.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindMethod {
			continue
		}
		owner, params, err := NewMethodParams(h4).Records()
		if err != nil {
			return out, err
		}
		for _, param := range params {
			out.Insert(model.FieldKey{Owner: owner, Key: param.Key}, param)
		}
	}
	return out, nil
}
