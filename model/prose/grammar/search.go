// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package grammar

import (
	"github.com/andreychh/tgen/model/prose"
)

// Search represents a scan of an inline run for the first position where a
// given [Rule] matches.
type Search[T any] struct {
	rule Rule[T]
}

// NewSearch constructs a Search over rule.
func NewSearch[T any](rule Rule[T]) Search[T] {
	return Search[T]{rule: rule}
}

// Find returns the value the rule decodes at the first position in inlines
// where it matches, and reports whether any position matched. The value is the
// zero value when the report is false.
func (s Search[T]) Find(inlines []prose.Inline) (T, bool) {
	for offset := range inlines {
		if value, ok := s.rule.Match(inlines[offset:]); ok {
			return value, true
		}
	}
	var zero T
	return zero, false
}
