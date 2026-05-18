// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/model/spec"
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

func (s Specification) Objects() iter.Seq[spec.Object] {
	return func(yield func(spec.Object) bool) {
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

func (s Specification) Methods() iter.Seq[spec.Method] {
	return func(yield func(spec.Method) bool) {
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

func (s Specification) DiscriminatedObjects() iter.Seq[spec.DiscriminatedObject] {
	return func(yield func(spec.DiscriminatedObject) bool) {
		seq := s.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewHeader(s.root, h4).Kind() == DefinitionKindDiscriminatedObject
			}).
			All()
		for h4 := range seq {
			if !yield(NewDiscriminatedObject(s.root, h4)) {
				break
			}
		}
	}
}

func (s Specification) DiscriminatedUnions() iter.Seq[spec.DiscriminatedUnion] {
	return func(yield func(spec.DiscriminatedUnion) bool) {
		seq := s.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewHeader(s.root, h4).Kind() == DefinitionKindDiscriminatedUnion
			}).
			All()
		for h4 := range seq {
			if !yield(NewDiscriminatedUnion(s.root, h4)) {
				break
			}
		}
	}
}

func (s Specification) Release() spec.Release {
	return NewRelease(s.root.Find("div#dev_page_content h4").At(0))
}
