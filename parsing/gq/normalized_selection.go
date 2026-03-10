// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

//nolint:ireturn // returns interface by design — Selection is a fluent interface
package gq

import (
	"iter"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type NormSelection struct {
	inner *goquery.Selection
}

func NewNormSelection(s *goquery.Selection) NormSelection {
	return NormSelection{inner: s}
}

func (s NormSelection) Text() string {
	return s.norm(s.inner.Text())
}

func (s NormSelection) Attr(name string) (string, bool) {
	val, exists := s.inner.Attr(name)
	if !exists {
		return "", false
	}
	return s.norm(val), true
}

func (s NormSelection) Find(selector string) Selection {
	return NewNormSelection(s.inner.Find(selector))
}

func (s NormSelection) Filter(selector string) Selection {
	return NewNormSelection(s.inner.Filter(selector))
}

func (s NormSelection) FilterFunc(f func(Selection) bool) Selection {
	return NewNormSelection(
		s.inner.FilterFunction(
			func(_ int, gs *goquery.Selection) bool {
				return f(NewNormSelection(gs))
			},
		),
	)
}

func (s NormSelection) Until(selector string) Selection {
	return NewNormSelection(s.inner.NextUntil(selector))
}

func (s NormSelection) Length() int {
	return s.inner.Length()
}

func (s NormSelection) IsEmpty() bool {
	return s.Length() == 0
}

func (s NormSelection) At(index int) Selection {
	if index < 0 {
		return NewNormSelection(s.inner.Slice(0, 0))
	}
	return NewNormSelection(s.inner.Eq(index))
}

func (s NormSelection) All() iter.Seq[Selection] {
	return func(yield func(Selection) bool) {
		for _, gs := range s.inner.EachIter() {
			if !yield(NewNormSelection(gs)) {
				return
			}
		}
	}
}

func (s NormSelection) norm(raw string) string {
	raw = strings.ReplaceAll(raw, "\n", " ")
	return strings.Join(strings.Fields(raw), " ")
}
