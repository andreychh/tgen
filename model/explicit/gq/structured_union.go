// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

// StructuredUnion represents a union whose variants are distinguished by their
// unique field sets. It navigates from the document root to locate each
// variant's full definition.
type StructuredUnion struct {
	root gq.Selection
	h4   gq.Selection
}

// NewStructuredUnion constructs a StructuredUnion from the document root and
// the h4 selection of the union definition.
func NewStructuredUnion(root, h4 gq.Selection) StructuredUnion {
	return StructuredUnion{root: root, h4: h4}
}

func (u StructuredUnion) Reference() model.Reference {
	return NewDefinitionReference(u.h4.Find("a.anchor"))
}

func (u StructuredUnion) Name() model.Name {
	return NewName(u.h4)
}

func (u StructuredUnion) Description() model.Description {
	return NewDefinitionDescription(u.h4)
}

// Variants returns the variant objects of this union. Each variant name is read
// from the union's <ul> list and resolved to a full object definition by
// navigating from the document root.
//
//nolint:varnamelen // <h4> is the standard HTML heading element name
func (u StructuredUnion) Variants() iter.Seq[explicit.Object] {
	return func(yield func(explicit.Object) bool) {
		for li := range u.h4.Until("h3, h4, hr").Find("ul li").All() {
			h4 := u.root.
				Find("div#dev_page_content h4").
				FilterFunc(func(s gq.Selection) bool {
					return s.Text() == li.Find("a").Text()
				}).
				At(0)
			if !yield(NewObject(h4)) {
				break
			}
		}
	}
}
