// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/model/pipeline/attached"
	"github.com/andreychh/tgen/model/pipeline/classified"
	"github.com/andreychh/tgen/model/pipeline/corrected"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/resolved"
	"github.com/andreychh/tgen/model/pipeline/separated"
	"github.com/andreychh/tgen/model/pipeline/typed"
	"github.com/andreychh/tgen/model/pipeline/unified"
)

// Pipeline represents the nanopass chain read end to end: every pass from the
// documentation page to the tables a target renders.
type Pipeline struct {
	doc *goquery.Document
}

// NewPipeline creates a Pipeline over a parsed documentation page.
func NewPipeline(doc *goquery.Document) Pipeline {
	return Pipeline{doc: doc}
}

// Specification returns the tables a target renders, as the last pass of the
// chain leaves them. It fails when any pass rejects what the pass before it
// produced.
func (p Pipeline) Specification() (separated.Specification, error) {
	page, err := parsed.NewPage(p.doc).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("parsing the page: %w", err)
	}
	discriminators, err := classified.NewPass(page).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("classifying discriminators: %w", err)
	}
	fields, err := unified.NewPass(discriminators).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("unifying fields and parameters: %w", err)
	}
	types, err := typed.NewPass(fields).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("typing fields: %w", err)
	}
	returns, err := resolved.NewPass(types).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("resolving returns: %w", err)
	}
	corrections, err := corrected.NewPass(returns).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("correcting definitions: %w", err)
	}
	forms, err := flattened.NewPass(corrections).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("flattening types: %w", err)
	}
	files, err := attached.NewPass(forms).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("attaching files: %w", err)
	}
	spec, err := separated.NewPass(files).Specification()
	if err != nil {
		return separated.Specification{}, fmt.Errorf("separating returns: %w", err)
	}
	return spec, nil
}
