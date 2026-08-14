// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package grammar

import (
	"regexp"

	"github.com/andreychh/tgen/model/prose"
)

// Marker represents a prose signal word — "Returns", "is returned", "must be"
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

// Capture represents a value a regular expression lifts out of a plain text
// run through its first capturing group, the form the documentation writes a
// fixed value in when it quotes one inside a sentence. A pattern that captures
// nothing matches nothing, so a group the rest of the pattern must not compete
// with is written as (?:…).
type Capture struct {
	pattern *regexp.Regexp
}

// NewCapture constructs a Capture from pattern.
func NewCapture(pattern *regexp.Regexp) Capture {
	return Capture{pattern: pattern}
}

// Matches returns the value the pattern's first capturing group lifts out of
// the inline, and reports whether the inline is a plain text run the pattern
// matches with that group filled. The value is empty when the report is false.
func (c Capture) Matches(inline prose.Inline) (string, bool) {
	text, ok := inline.(prose.Text)
	if !ok || text.Style() != prose.StylePlain {
		return "", false
	}
	match := c.pattern.FindStringSubmatch(text.Content())
	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}

// Italic represents an emphasized text run, the form the documentation writes
// a value in when it sets one apart from the sentence around it.
type Italic struct{}

// NewItalic constructs an Italic.
func NewItalic() Italic {
	return Italic{}
}

// Matches returns the content of the inline verbatim, and reports whether the
// inline is an italic text run. The content is empty when the report is false.
func (Italic) Matches(inline prose.Inline) (string, bool) {
	text, ok := inline.(prose.Text)
	if !ok || text.Style() != prose.StyleItalic {
		return "", false
	}
	return text.Content(), true
}

// Anchor represents an in-page link, the form the documentation identifies one
// of its own sections with.
type Anchor struct{}

// NewAnchor constructs an Anchor.
func NewAnchor() Anchor {
	return Anchor{}
}

// Matches returns the fragment the inline addresses, without its leading "#",
// and reports whether the inline is a link to an in-page anchor. The fragment
// is empty when the report is false.
func (Anchor) Matches(inline prose.Inline) (string, bool) {
	link, ok := inline.(prose.Link)
	if !ok {
		return "", false
	}
	return link.Anchor()
}
