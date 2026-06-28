// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

import "strings"

// Style is the emphasis applied to a text run.
type Style int

const (
	// StylePlain marks a run with no emphasis.
	StylePlain Style = iota
	// StyleItalic marks an emphasized run.
	StyleItalic
	// StyleBold marks a strongly emphasized run.
	StyleBold
	// StyleCode marks a verbatim monowidth run.
	StyleCode
)

// Text represents a run of text rendered with a single style.
type Text struct {
	content string
	style   Style
}

// NewText constructs a styled text run from its content and style.
func NewText(content string, style Style) Text {
	return Text{content: content, style: style}
}

// Content returns the text of the run.
func (t Text) Content() string {
	return t.content
}

// Style returns the emphasis applied to the run.
func (t Text) Style() Style {
	return t.style
}

// HasPrefix reports whether the run's content begins with prefix.
func (t Text) HasPrefix(prefix string) bool {
	return strings.HasPrefix(t.content, prefix)
}

// TrimPrefix returns the run with prefix removed from the front of its content,
// keeping its style. The run is returned unchanged when its content lacks the
// prefix.
func (t Text) TrimPrefix(prefix string) Text {
	return Text{content: strings.TrimPrefix(t.content, prefix), style: t.style}
}

// Equals reports whether other is a Text with the same content and style.
func (t Text) Equals(other Inline) bool {
	o, ok := other.(Text)
	return ok && t.content == o.content && t.style == o.style
}

func (Text) isInline() {}

// Link represents a styled text run that addresses a target URL or anchor.
type Link struct {
	content string
	style   Style
	href    string
}

// NewLink constructs a link to href from a styled run of content.
func NewLink(content string, style Style, href string) Link {
	return Link{content: content, style: style, href: href}
}

// Content returns the text of the link.
func (l Link) Content() string {
	return l.content
}

// Style returns the emphasis applied to the link.
func (l Link) Style() Style {
	return l.style
}

// Href returns the URL or anchor the link addresses.
func (l Link) Href() string {
	return l.href
}

// Telegram Bot API 10.1 uses three href forms across the documentation prose:
//   - "#section"  — an anchor into the same page; the only form inside type cells
//   - "https://…" — an absolute URL to an external resource (descriptions only)
//   - "/path"     — a path relative to the documentation site root (descriptions only)
//
// Only the anchor form denotes a documented type, so that is the one distinction
// Anchor draws; the others may earn their own modeling once descriptions are
// rendered into absolute links.

// Anchor returns the fragment the link targets and reports whether the link is
// an anchor — an in-page reference of the form "#section". The fragment is
// returned without its leading "#", and is empty when the link is not an anchor.
func (l Link) Anchor() (string, bool) {
	if !strings.HasPrefix(l.href, "#") {
		return "", false
	}
	return strings.TrimPrefix(l.href, "#"), true
}

// Equals reports whether other is a Link with the same content, style, and href.
func (l Link) Equals(other Inline) bool {
	o, ok := other.(Link)
	return ok && l.content == o.content && l.style == o.style && l.href == o.href
}

func (Link) isInline() {}

// LineBreak represents a forced line break within a block.
type LineBreak struct{}

// NewLineBreak constructs a line break.
func NewLineBreak() LineBreak {
	return LineBreak{}
}

// Equals reports whether other is also a LineBreak.
func (LineBreak) Equals(other Inline) bool {
	_, ok := other.(LineBreak)
	return ok
}

func (LineBreak) isInline() {}
