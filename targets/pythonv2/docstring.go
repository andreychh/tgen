// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"strings"

	"github.com/andreychh/tgen/model/prose"
	"github.com/mitchellh/go-wordwrap"
)

// width is the room a line of flowing text has, indentation counted in. It is
// the limit PEP 8 puts on a docstring, which is narrower than the one it puts
// on code.
const width = 72

// quotes opens and closes a docstring.
const quotes = `"""`

// Docstring represents a prose passage rendered as a Python docstring: a
// paragraph per block, a list as one dashed line per item, and the whole
// enclosed in triple quotes. The first line carries no indentation, since
// whatever declares the documented name has already written it.
type Docstring struct {
	passage prose.Passage
	indent  int
}

// NewDocstring creates a Docstring rendering a passage at an indentation depth.
func NewDocstring(passage prose.Passage, indent int) Docstring {
	return Docstring{passage: passage, indent: indent}
}

// NewClassDocstring creates a Docstring for the body of a class.
func NewClassDocstring(passage prose.Passage) Docstring {
	return NewDocstring(passage, 4)
}

// NewStatementDocstring creates a Docstring for a name declared by a statement
// at the margin of the module, which nothing indents.
func NewStatementDocstring(passage prose.Passage) Docstring {
	return NewDocstring(passage, 0)
}

// NewFieldDocstring creates a Docstring for a name declared inside a class. A
// field is described by a table cell, which holds inline prose only, so its one
// phrase becomes the single paragraph of a passage.
func NewFieldDocstring(phrase prose.Phrase) Docstring {
	return NewDocstring(prose.NewPassage(prose.NewParagraph(phrase.Inlines()...)), 4)
}

// Value returns the docstring, empty when the passage writes no prose. A block
// that writes nothing takes no line and earns no blank line beside it, so an
// empty paragraph leaves no trace rather than an empty docstring.
func (d Docstring) Value() string {
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
	return d.quote(lines)
}

// block returns the lines one block occupies.
func (d Docstring) block(block prose.Block) []string {
	switch block := block.(type) {
	case prose.Paragraph:
		return fold(join(block.Inlines()), width-d.indent, "", "")
	case prose.List:
		return d.list(block)
	default:
		return nil
	}
}

// list returns the lines a list occupies, each item written as a dashed line.
// Python states no form for a list inside a docstring, so what is written is
// what a reader of the source reads as one.
func (d Docstring) list(list prose.List) []string {
	lines := make([]string, 0, len(list.Items()))
	for _, item := range list.Items() {
		lines = append(lines, fold(join(item.Inlines()), width-d.indent-2, "- ", "  ")...)
	}
	return lines
}

// quote returns the lines enclosed in triple quotes, every line after the first
// indented to the depth the declaration sits at. One line closes on the line it
// opened, unless the prose ends in a quote of its own, which would close the
// docstring with a run of four.
func (d Docstring) quote(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	if len(lines) == 1 && !strings.HasSuffix(lines[0], `"`) {
		return quotes + lines[0] + quotes
	}
	pad := strings.Repeat(" ", d.indent)
	out := make([]string, 0, len(lines)+2)
	out = append(out, quotes+lines[0])
	for _, line := range lines[1:] {
		out = append(out, strings.TrimRight(pad+line, " "))
	}
	return strings.Join(append(out, pad+quotes), "\n")
}

// join returns the plain text of inline content. A link contributes its text
// alone: the anchor it addresses is not yet resolved to the name a docstring
// would have to spell.
func join(inlines []prose.Inline) string {
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

// fold returns content folded to the given width, opening with first and
// continuing with rest. A forced line break in the content starts a new line of
// its own. Content with nothing to read folds to no lines at all.
func fold(content string, width int, first, rest string) []string {
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
