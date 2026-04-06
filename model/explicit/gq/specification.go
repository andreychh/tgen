// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

type Specification struct {
	root gq.Selection
}

func NewSpecification(root gq.Selection) Specification {
	return Specification{root: root}
}

// NewSpecificationFromDocument creates a Specification from a parsed goquery document.
func NewSpecificationFromDocument(doc *goquery.Document) Specification {
	return NewSpecification(gq.NewNormSelection(doc.Selection))
}

func (s Specification) Objects() iter.Seq[explicit.Object] {
	return func(yield func(explicit.Object) bool) {
		seq := s.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewHeader(s.root, h4).Kind() == DefinitionKindObject
			}).
			All()
		for h4 := range seq {
			if !yield(NewObject(s.root, h4)) {
				break
			}
		}
	}
}

func (s Specification) Methods() iter.Seq[explicit.Method] {
	return func(yield func(explicit.Method) bool) {
		seq := s.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewHeader(s.root, h4).Kind() == DefinitionKindMethod
			}).
			All()
		for h4 := range seq {
			if !yield(NewMethod(s.root, h4)) {
				break
			}
		}
	}
}

func (s Specification) Unions() explicit.Unions {
	return NewUnions(s.root)
}

func (s Specification) Release() explicit.Release {
	return NewRelease(s.root.Find("div#dev_page_content h4").At(0))
}
