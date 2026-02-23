// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/dom"
)

// RawSpecification implements the Specification interface by parsing a raw HTML
// DOM tree.
//
// It serves as the main entry point for extracting API entities directly from
// the Telegram Bot API documentation page.
type RawSpecification struct {
	selection dom.Selection
}

// NewRawSpecification creates a new RawSpecification instance.
//
// The provided root represents the parsed HTML document or the specific
// container element holding the API documentation content.
func NewRawSpecification(root dom.Selection) RawSpecification {
	return RawSpecification{selection: root}
}

// Unions returns an iterator over all polymorphic sum-types (Union) defined in
// the specification.
// TODO: consider using `div#dev_page_content h4` selector.
func (s RawSpecification) Unions() iter.Seq[Union] {
	return func(yield func(Union) bool) {
		seq := s.selection.Find("h4").FilterFunc(
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

// Objects returns an iterator over all standard object defined in the
// specification.
// TODO: consider using `div#dev_page_content h4` selector.
func (s RawSpecification) Objects() iter.Seq[Object] {
	return func(yield func(Object) bool) {
		seq := s.selection.Find("h4").FilterFunc(
			func(s dom.Selection) bool {
				return NewAnchor(s).Kind() == KindObject
			},
		).All()
		for _, h4 := range seq {
			if !yield(NewRawObject(h4)) {
				break
			}
		}
	}
}
