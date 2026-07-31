// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package ir reads the normalized tables of the last pipeline stage as the
// records a target renders: it joins a definition with what it owns, binds
// every type to the definition it names, and puts both in the order their
// source gave them. It is the pipeline's exit — nothing rewrites the
// specification after it, and no target reaches past it to a table.
//
// Records are read as one sequence rather than one per kind. The page numbers
// its headings in a single run, so an object, a union and a method stand beside
// each other in it, and a target walking the sequence writes its file the way
// the page reads. A definition tgen introduces has no place on the page and
// follows every definition that has one.
package ir

import (
	"github.com/andreychh/tgen/model/pipeline/separated"
)

// Specification represents the pipeline's database read as records to render.
type Specification struct {
	db separated.Specification
}

// NewSpecification constructs a Specification over a separated specification.
func NewSpecification(db separated.Specification) Specification {
	return Specification{db: db}
}

// Definitions returns every definition the specification names, each joined
// with what it owns, ordered by the position its source gave it. It fails when
// a definition cannot be read as the record of its kind.
func (s Specification) Definitions() ([]Definition, error) {
	out := make([]Definition, 0)
	for _, definition := range NewOrdered(s.db.Definitions).Value() {
		record, err := NewReading(s.db, definition).Value()
		if err != nil {
			return nil, err
		}
		out = append(out, record)
	}
	return out, nil
}
