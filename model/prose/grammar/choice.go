// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package grammar

import (
	"github.com/andreychh/tgen/model/prose"
)

// Choice is a [Rule] standing for a set of alternatives, tried in the order
// they were given.
type Choice[T any] struct {
	rules []Rule[T]
}

// NewChoice constructs a Choice from its alternatives.
func NewChoice[T any](rules ...Rule[T]) Choice[T] {
	return Choice[T]{rules: rules}
}

// Match implements [Rule]. It returns what the first alternative to match
// decodes, leaving the alternatives after it untried.
func (c Choice[T]) Match(inlines []prose.Inline) (T, bool) {
	for _, rule := range c.rules {
		if value, ok := rule.Match(inlines); ok {
			return value, true
		}
	}
	var zero T
	return zero, false
}
