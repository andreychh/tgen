// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/types/v2"
)

// Search represents a scan of an inline run for the first position where a
// given [Rule] matches.
type Search struct {
	rule Rule
}

// NewSearch constructs a Search over rule.
func NewSearch(rule Rule) Search {
	return Search{rule: rule}
}

// Find returns the type expression decoded at the first position in inlines
// where the rule matches, or false when no position matches.
func (s Search) Find(inlines []prose.Inline) (types.Expression, bool) {
	for offset := range inlines {
		if expr, ok := s.rule.Match(inlines[offset:]); ok {
			return expr, true
		}
	}
	return nil, false
}
