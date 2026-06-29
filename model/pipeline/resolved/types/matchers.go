// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package types decodes the prose of a method's return clause into a type
// expression.
//
// [ReturnType] is the entry point: wrap a method's description prose and call
// its Value method to obtain the return type.
//
// Decoding is driven by [Rule] implementations — [ReturnsRule], [ReturnedRule],
// [ArrayRule], and [UnionRule] — each recognizing one structural form of the
// clause. [ProductionRule] tries them in priority order; [Search] applies the
// combined rule at every position in the paragraph until one matches.
package types

import (
	"regexp"
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/types/v2"
)

// Marker represents a prose signal word — "Returns", "is returned", "Array of"
// — recognized by a regular expression against plain text runs.
type Marker struct {
	pattern *regexp.Regexp
}

// NewMarker constructs a Marker from pattern.
func NewMarker(pattern *regexp.Regexp) Marker {
	return Marker{pattern: pattern}
}

// Matches reports whether the inline is a plain text run whose content
// satisfies the marker's regular expression.
func (m Marker) Matches(inline prose.Inline) bool {
	text, ok := inline.(prose.Text)
	return ok && text.Style() == prose.StylePlain && m.pattern.MatchString(text.Content())
}

// Named represents an in-page anchor link, the form a return clause uses to
// identify a documented type.
type Named struct{}

// NewNamed constructs a Named.
func NewNamed() Named {
	return Named{}
}

// Matches reports whether the inline is a link to an in-page anchor and returns
// the reference it addresses. The reference is empty when the report is false.
func (Named) Matches(inline prose.Inline) (model.Reference, bool) {
	link, ok := inline.(prose.Link)
	if !ok {
		return "", false
	}
	target, ok := link.Anchor()
	if !ok {
		return "", false
	}
	return model.Reference(target), true
}

// Primitive represents the italic keyword form that a return clause uses for
// built-in types such as True, String, and Float.
type Primitive struct {
	primitives types.Primitives
}

// NewPrimitive constructs a Primitive over the default primitive vocabulary.
func NewPrimitive() Primitive {
	return Primitive{primitives: types.NewPrimitives()}
}

// Matches reports whether the inline is an italic text run naming a built-in
// type and returns the kind it names. The kind is empty when the report is
// false.
func (p Primitive) Matches(inline prose.Inline) (types.PrimitiveKind, bool) {
	text, ok := inline.(prose.Text)
	if !ok || text.Style() != prose.StyleItalic {
		return "", false
	}
	return p.primitives.Kind(strings.TrimSpace(text.Content()))
}
