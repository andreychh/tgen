// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

//nolint:ireturn // returns interface by design — Selection is a fluent interface
package gq

import (
	"iter"

	"github.com/PuerkitoBio/goquery"
)

type RawSelection struct {
	inner *goquery.Selection
}

func NewRawSelection(s *goquery.Selection) RawSelection {
	return RawSelection{inner: s}
}

func (s RawSelection) Text() string {
	return s.inner.Text()
}

func (s RawSelection) Attr(name string) (string, bool) {
	return s.inner.Attr(name)
}

func (s RawSelection) Find(selector string) Selection {
	return NewRawSelection(s.inner.Find(selector))
}

func (s RawSelection) Filter(selector string) Selection {
	return NewRawSelection(s.inner.Filter(selector))
}

func (s RawSelection) FilterFunc(f func(Selection) bool) Selection {
	return NewRawSelection(s.inner.FilterFunction(func(_ int, gs *goquery.Selection) bool {
		return f(NewRawSelection(gs))
	}))
}

func (s RawSelection) Until(selector string) Selection {
	return NewRawSelection(s.inner.NextUntil(selector))
}

func (s RawSelection) IsEmpty() bool {
	return s.inner.Length() == 0
}

func (s RawSelection) At(index int) Selection {
	if index < 0 {
		return NewRawSelection(s.inner.Slice(0, 0))
	}
	return NewRawSelection(s.inner.Eq(index))
}

func (s RawSelection) All() iter.Seq[Selection] {
	return func(yield func(Selection) bool) {
		for _, gqs := range s.inner.EachIter() {
			if !yield(NewRawSelection(gqs)) {
				return
			}
		}
	}
}
