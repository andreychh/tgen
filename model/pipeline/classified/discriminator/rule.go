// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package discriminator

import (
	"regexp"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/prose/grammar"
)

var (
	alwaysSignal = regexp.MustCompile(`(?i)\balways\s*“([^”]+)”\s*$`)
	mustBeSignal = regexp.MustCompile(`(?i)\bmust be\s*$`)
)

// Rule is a pattern that recognizes one structural form of a fixed
// discriminator value inside a field's description, decoding it into the value
// it names.
type Rule = grammar.Rule[model.DiscriminatorValue]

// AlwaysRule is a [Rule] that recognizes "…, always “value”", a fixed value
// quoted inside a single plain text run.
type AlwaysRule struct{}

// NewAlwaysRule constructs an AlwaysRule.
func NewAlwaysRule() AlwaysRule {
	return AlwaysRule{}
}

// Match implements [Rule].
func (AlwaysRule) Match(inlines []prose.Inline) (model.DiscriminatorValue, bool) {
	if len(inlines) < 1 {
		return "", false
	}
	value, ok := grammar.NewCapture(alwaysSignal).Matches(inlines[0])
	if !ok {
		return "", false
	}
	return model.DiscriminatorValue(value), true
}

// MustBeRule is a [Rule] that recognizes "…, must be <value>", a fixed value
// carried by a trailing italic run.
type MustBeRule struct{}

// NewMustBeRule constructs a MustBeRule.
func NewMustBeRule() MustBeRule {
	return MustBeRule{}
}

// Match implements [Rule].
func (MustBeRule) Match(inlines []prose.Inline) (model.DiscriminatorValue, bool) {
	if len(inlines) < 2 {
		return "", false
	}
	if !grammar.NewMarker(mustBeSignal).Matches(inlines[0]) {
		return "", false
	}
	value, ok := grammar.NewItalic().Matches(inlines[1])
	if !ok {
		return "", false
	}
	return model.DiscriminatorValue(value), true
}
