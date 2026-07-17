// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package discriminator

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
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
func (p ProductionRule) Match(inlines []prose.Inline) (model.DiscriminatorValue, bool) {
	for _, rule := range p.rules {
		if value, ok := rule.Match(inlines); ok {
			return value, true
		}
	}
	return "", false
}
