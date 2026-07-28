// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

// ProductionRule is a [Rule] composed of other rules tried in priority order.
type ProductionRule struct {
	rules []Rule
}

// NewProductionRule constructs a ProductionRule from its rules, in priority
// order.
func NewProductionRule(rules ...Rule) ProductionRule {
	return ProductionRule{rules: rules}
}

// Match implements [Rule].
func (p ProductionRule) Match(inlines []prose.Inline) (typeexpr.Expression, bool) {
	for _, rule := range p.rules {
		if expr, ok := rule.Match(inlines); ok {
			return expr, true
		}
	}
	return nil, false
}
