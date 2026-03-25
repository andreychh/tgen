// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type Specification struct {
	selection gq.Selection
}

func NewSpecification(root gq.Selection) Specification {
	return Specification{selection: root}
}

func (s Specification) Objects() iter.Seq[Object] {
	return func(yield func(Object) bool) {
		seq := s.selection.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewDefinitionHeader(s.selection, h4).Kind() == KindObject
			}).
			All()
		for h4 := range seq {
			if !yield(NewObject(h4)) {
				break
			}
		}
	}
}

func (s Specification) VariantObjects() iter.Seq[Object] {
	return func(yield func(Object) bool) {
		seq := s.selection.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewDefinitionHeader(s.selection, h4).Kind() == KindVariantObject
			}).
			All()
		for h4 := range seq {
			if !yield(NewObject(h4)) {
				break
			}
		}
	}
}

func (s Specification) Methods() iter.Seq[Method] {
	return func(yield func(method Method) bool) {
		seq := s.selection.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewDefinitionHeader(s.selection, h4).Kind() == KindMethod
			}).
			All()
		for h4 := range seq {
			if !yield(NewMethod(h4)) {
				break
			}
		}
	}
}

func (s Specification) Unions() iter.Seq[Union] {
	return func(yield func(Union) bool) {
		seq := s.selection.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewDefinitionHeader(s.selection, h4).Kind() == KindUnion
			}).
			All()
		for h4 := range seq {
			if !yield(NewUnion(h4)) {
				break
			}
		}
	}
}

func (s Specification) DiscriminatedUnions() iter.Seq[DiscriminatedUnion] {
	return func(yield func(DiscriminatedUnion) bool) {
		seq := s.selection.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewDefinitionHeader(s.selection, h4).Kind() == KindDiscriminatedUnion
			}).
			All()
		for h4 := range seq {
			if !yield(NewDiscriminatedUnion(s.selection, h4)) {
				break
			}
		}
	}
}

func (s Specification) Release() Release {
	return NewRelease(s.selection.Find("div#dev_page_content h4").At(0))
}
