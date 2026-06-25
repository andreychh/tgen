// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

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

func (Link) isInline() {}

// LineBreak represents a forced line break within a block.
type LineBreak struct{}

// NewLineBreak constructs a line break.
func NewLineBreak() LineBreak {
	return LineBreak{}
}

func (LineBreak) isInline() {}
