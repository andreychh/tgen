// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/dom"
)

// RawDocument implements the Document interface using a DOM selection as its
// data source.
type RawDocument struct {
	selection dom.Selection
}

// NewRawDocument creates a new RawDocument from the provided selection.
func NewRawDocument(s dom.Selection) RawDocument {
	return RawDocument{selection: s}
}

// Unions returns an iterator that filters and parses union types from the
// document's HTML structure.
func (u RawDocument) Unions() iter.Seq[Union] {
	return func(yield func(Union) bool) {
		seq := u.selection.Find("h4").FilterFunc(
			func(s dom.Selection) bool {
				return NewAnchor(s).Kind() == KindUnion
			},
		).All()
		for _, h4 := range seq {
			if !yield(NewDefaultRawUnion(h4)) {
				break
			}
		}
	}
}
