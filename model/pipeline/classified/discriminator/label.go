// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package discriminator decodes the prose of a field's description into the
// fixed discriminator value it may carry.
//
// [Label] is the entry point: wrap a field's description and call its Value
// method. Decoding is driven by [Rule] implementations — [AlwaysRule] and
// [MustBeRule] — each recognizing one structural form a discriminator value
// takes; a choice from [grammar] takes the first of them to match. A
// description matching neither rule is not a label, which Value reports rather
// than treats as failure.
package discriminator

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/prose/grammar"
)

// Label is a field's description ready to be decoded into the fixed
// discriminator value it may carry.
type Label struct {
	description prose.Phrase
}

// NewLabel constructs a Label over a field's description.
func NewLabel(description prose.Phrase) Label {
	return Label{description: description}
}

// Value returns the discriminator value decoded from the description, and
// reports whether the description carries one.
func (l Label) Value() (model.DiscriminatorValue, bool) {
	return l.rules().Match(l.description.Inlines())
}

// rules assembles the discriminator alternatives. Both anchor their form to
// the end of the description — one to a quoted value, the other to "must be" —
// so no description reaches both, and the order they stand in decides nothing.
func (l Label) rules() grammar.Choice[model.DiscriminatorValue] {
	return grammar.NewChoice(NewAlwaysRule(), NewMustBeRule())
}
