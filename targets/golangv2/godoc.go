// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"strings"

	"github.com/andreychh/tgen/model/prose"
	"github.com/mitchellh/go-wordwrap"
)

// width is the room left for doc text once the comment marker is written.
const width = 77

// Godoc represents a prose passage rendered as a Go doc comment: a paragraph
// per block, lists in the bullet form gofmt reflows, and one comment marker per
// line. The first line carries no indentation, since whatever declares the
// commented name has already written it.
type Godoc struct {
	passage prose.Passage
	indent  int
}

// NewGodoc creates a Godoc rendering a passage at an indentation depth.
func NewGodoc(passage prose.Passage, indent int) Godoc {
	return Godoc{passage: passage, indent: indent}
}

// NewTypeGodoc creates a Godoc for a name declared at file scope.
func NewTypeGodoc(passage prose.Passage) Godoc {
	return NewGodoc(passage, 0)
}

// NewFieldGodoc creates a Godoc for a name declared inside a struct. A field
// is described by a table cell, which holds inline prose only, so its one
// phrase becomes the single paragraph of a passage.
func NewFieldGodoc(phrase prose.Phrase) Godoc {
	return NewGodoc(prose.NewPassage(prose.NewParagraph(phrase.Inlines()...)), 1)
}

// Value returns the doc comment, empty when the passage writes no prose. A
// block that writes nothing takes no line and earns no blank line beside it,
// so an empty paragraph leaves no trace rather than an empty comment.
func (d Godoc) Value() string {
	lines := make([]string, 0)
	for _, block := range d.passage.Blocks() {
		written := d.block(block)
		if len(written) == 0 {
			continue
		}
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, written...)
	}
	return d.comment(lines)
}

// block returns the lines one block occupies.
func (d Godoc) block(block prose.Block) []string {
	switch block := block.(type) {
	case prose.Paragraph:
		return wrap(text(block.Inlines()), width, "", "")
	case prose.List:
		return d.list(block)
	default:
		return nil
	}
}

// list returns the lines a list occupies, each item written in the bullet form
// a Go doc comment recognises.
func (d Godoc) list(list prose.List) []string {
	lines := make([]string, 0, len(list.Items()))
	for _, item := range list.Items() {
		lines = append(lines, wrap(text(item.Inlines()), width-4, "  - ", "    ")...)
	}
	return lines
}

// comment returns the lines written as a doc comment, every line after the
// first indented to the depth the declaration sits at.
func (d Godoc) comment(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	pad := strings.Repeat("\t", d.indent)
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		out = append(out, strings.TrimRight(marker(pad, i == 0)+line, " "))
	}
	return strings.Join(out, "\n")
}

// marker returns the comment marker a line opens with.
func marker(pad string, first bool) string {
	if first {
		return "// "
	}
	return pad + "// "
}

// text returns the plain text of inline content. A link contributes its text
// alone: the anchor it addresses is not yet resolved to the name a Go doc link
// would have to spell.
func text(inlines []prose.Inline) string {
	var out strings.Builder
	for _, inline := range inlines {
		switch inline := inline.(type) {
		case prose.Text:
			out.WriteString(inline.Content())
		case prose.Link:
			out.WriteString(inline.Content())
		case prose.LineBreak:
			out.WriteString("\n")
		}
	}
	return out.String()
}

// wrap returns content folded to the given width, opening with first and
// continuing with rest. A forced line break in the content starts a new line of
// its own. Content with nothing to read folds to no lines at all.
func wrap(content string, width int, first, rest string) []string {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	out := make([]string, 0)
	for segment := range strings.SplitSeq(content, "\n") {
		folded := wordwrap.WrapString(segment, uint(width))
		for line := range strings.SplitSeq(folded, "\n") {
			out = append(out, rest+line)
		}
	}
	out[0] = first + strings.TrimPrefix(out[0], rest)
	return out
}
