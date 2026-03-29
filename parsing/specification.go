// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQSpecification struct {
	root gq.Selection
}

func NewGQSpecification(root gq.Selection) GQSpecification {
	return GQSpecification{root: root}
}

func (s GQSpecification) Objects() iter.Seq[Object] {
	return func(yield func(Object) bool) {
		seq := s.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewGQHeader(s.root, h4).Kind() == DefinitionKindObject
			}).
			All()
		for h4 := range seq {
			if !yield(NewGQObject(h4)) {
				break
			}
		}
	}
}

func (s GQSpecification) Methods() iter.Seq[Method] {
	return func(yield func(Method) bool) {
		seq := s.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewGQHeader(s.root, h4).Kind() == DefinitionKindMethod
			}).
			All()
		for h4 := range seq {
			if !yield(NewGQMethod(h4)) {
				break
			}
		}
	}
}

func (s GQSpecification) Unions() Unions {
	return NewGQUnions(s.root)
}

func (s GQSpecification) Release() Release {
	return NewGQRelease(s.root.Find("div#dev_page_content h4").At(0))
}
