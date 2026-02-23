// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

//nolint:ireturn // Fluent API design requires returning the interface
package dom

import (
	"iter"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLSelection is the implementation of the Selection interface based on the
// [goquery] package.
type HTMLSelection struct {
	inner *goquery.Selection
}

// NewHTMLSelection creates a new HTMLSelection from a goquery.Selection.
func NewHTMLSelection(s *goquery.Selection) HTMLSelection {
	return HTMLSelection{inner: s}
}

// Text returns the normalized text of the selection.
func (s HTMLSelection) Text() string {
	return s.normalize(s.inner.Text())
}

// Attr returns the normalized value of the specified attribute.
func (s HTMLSelection) Attr(name string) (value string, exists bool) {
	value, exists = s.inner.Attr(name)
	if !exists {
		return "", false
	}
	return s.normalize(value), true
}

// First returns the first element of the selection.
func (s HTMLSelection) First() Selection {
	return NewHTMLSelection(s.inner.First())
}

// Find looks for descendants matching the selector.
func (s HTMLSelection) Find(selector string) Selection {
	return NewHTMLSelection(s.inner.Find(selector))
}

// Filter reduces the set to elements matching the selector.
func (s HTMLSelection) Filter(selector string) Selection {
	return NewHTMLSelection(s.inner.Filter(selector))
}

// FilterFunc reduces the set based on a predicate function.
func (s HTMLSelection) FilterFunc(f func(Selection) bool) Selection {
	return NewHTMLSelection(
		s.inner.FilterFunction(
			func(_ int, gs *goquery.Selection) bool {
				return f(NewHTMLSelection(gs))
			},
		),
	)
}

// NextAllFiltered returns all following siblings matching the selector.
func (s HTMLSelection) NextAllFiltered(selector string) Selection {
	return NewHTMLSelection(s.inner.NextAllFiltered(selector))
}

// NextUntil gets following siblings up to but not including the selector.
func (s HTMLSelection) NextUntil(selector string) Selection {
	return NewHTMLSelection(s.inner.NextUntil(selector))
}

// Length returns the number of elements in the selection.
func (s HTMLSelection) Length() int {
	return s.inner.Length()
}

// IsEmpty reports whether the selection contains no elements.
func (s HTMLSelection) IsEmpty() bool {
	return s.Length() == 0
}

// At returns the element at the specified index. Unlike jQuery/goquery,
// negative indices are not supported and will return an empty selection.
func (s HTMLSelection) At(index int) Selection {
	if index < 0 {
		return NewHTMLSelection(s.inner.Slice(0, 0))
	}
	return NewHTMLSelection(s.inner.Eq(index))
}

// All returns an iterator over the selection.
func (s HTMLSelection) All() iter.Seq2[int, Selection] {
	return func(yield func(int, Selection) bool) {
		for i, gs := range s.inner.EachIter() {
			if !yield(i, NewHTMLSelection(gs)) {
				return
			}
		}
	}
}

// normalize collapses whitespace sequences and trims the result.
func (s HTMLSelection) normalize(raw string) string {
	return strings.Join(strings.Fields(raw), " ")
}
