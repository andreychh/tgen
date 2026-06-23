// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

// Text represents a run of plain text.
type Text struct {
	content string
}

// NewText constructs a plain-text run from its content.
func NewText(content string) Text {
	return Text{content: content}
}

// Content returns the text of the run.
func (t Text) Content() string {
	return t.content
}

func (Text) isInline() {}

// Bold represents inline content rendered with strong emphasis.
type Bold struct {
	inlines []Inline
}

// NewBold constructs bold content from inline children.
func NewBold(inlines ...Inline) Bold {
	return Bold{inlines: inlines}
}

// Inlines returns the inline children of the bold run.
func (b Bold) Inlines() []Inline {
	return b.inlines
}

func (Bold) isInline() {}

// Italic represents inline content rendered with emphasis.
type Italic struct {
	inlines []Inline
}

// NewItalic constructs emphasized content from inline children.
func NewItalic(inlines ...Inline) Italic {
	return Italic{inlines: inlines}
}

// Inlines returns the inline children of the emphasized run.
func (i Italic) Inlines() []Inline {
	return i.inlines
}

func (Italic) isInline() {}

// Code represents a verbatim monowidth span.
type Code struct {
	content string
}

// NewCode constructs a verbatim span from its content.
func NewCode(content string) Code {
	return Code{content: content}
}

// Content returns the verbatim text of the span.
func (c Code) Content() string {
	return c.content
}

func (Code) isInline() {}

// Link represents inline content addressing a target URL or anchor.
type Link struct {
	target  string
	inlines []Inline
}

// NewLink constructs a link to target from inline children.
func NewLink(target string, inlines ...Inline) Link {
	return Link{target: target, inlines: inlines}
}

// Target returns the URL or anchor the link addresses.
func (l Link) Target() string {
	return l.target
}

// Inlines returns the inline children of the link.
func (l Link) Inlines() []Inline {
	return l.inlines
}

func (Link) isInline() {}

// LineBreak represents a forced line break within a block.
type LineBreak struct{}

// NewLineBreak constructs a line break.
func NewLineBreak() LineBreak {
	return LineBreak{}
}

func (LineBreak) isInline() {}
