// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package discriminator

import (
	"regexp"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
)

var (
	alwaysSignal = regexp.MustCompile(`(?i)\balways\s*“([^”]+)”\s*$`)
	mustBeSignal = regexp.MustCompile(`(?i)\bmust be\s*$`)
)

// Rule is a pattern that recognizes one structural form of a fixed
// discriminator value inside a field's description, decoding it into the value
// it names.
type Rule interface {
	// Match reports whether the rule's form matches the front of inlines and
	// returns the discriminator value it names.
	Match(inlines []prose.Inline) (model.DiscriminatorValue, bool)
}

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
	text, ok := inlines[0].(prose.Text)
	if !ok || text.Style() != prose.StylePlain {
		return "", false
	}
	match := alwaysSignal.FindStringSubmatch(text.Content())
	if match == nil {
		return "", false
	}
	return model.DiscriminatorValue(match[1]), true
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
	lead, ok := inlines[0].(prose.Text)
	if !ok || lead.Style() != prose.StylePlain || !mustBeSignal.MatchString(lead.Content()) {
		return "", false
	}
	value, ok := inlines[1].(prose.Text)
	if !ok || value.Style() != prose.StyleItalic {
		return "", false
	}
	return model.DiscriminatorValue(value.Content()), true
}
