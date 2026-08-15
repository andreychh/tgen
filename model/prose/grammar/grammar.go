// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package grammar supplies the pieces a pass builds its prose grammar out of.
//
// The rules belong to the pass, along with the vocabulary they name values in;
// this package holds what tries them. [Rule] is the shape every rule has,
// [Choice] takes the first of a set to match, and [Search] applies one at every
// position of a run. The matchers — [Marker], [Capture], [Italic], [Anchor] —
// read a single inline, which is what a rule composes its form out of.
//
// The value a rule decodes is what the passes differ in, and it is the type
// parameter every rule here carries.
package grammar

import (
	"github.com/andreychh/tgen/model/prose"
)

// Rule is a pattern that recognizes one structural form at the front of a run
// of inlines, decoding it into the value that form names.
type Rule[T any] interface {
	// Match reports whether the rule's form begins at the front of inlines and
	// returns the value it names. It reports false when the form does not start
	// at inlines[0], leaving the caller to advance or try another rule. The
	// value is the zero value when the report is false.
	Match(inlines []prose.Inline) (T, bool)
}
