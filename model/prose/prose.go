// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package prose models Telegram documentation text as a target-neutral
// structure of block and inline nodes, decoupled from both the source HTML and
// any output target. The decoder maps HTML into this model; each target renders
// it into its own doc-comment syntax. Section prose, with its paragraphs and
// lists, is a [Passage] of blocks; a table cell, whose text is inline only, is
// a [Phrase] of runs. Either way the model names the text, not its role — a
// field type as much as a field description.
//
// The model covers only the prose found inside entity sections — objects,
// methods, and unions. Markup that appears elsewhere on the page (the intro,
// the changelog, standalone guide sections) is out of scope, so constructs that
// occur only there have no node: there is no image, because real images live in
// the intro; no code block, because sections never fence one; and no ordered
// list, because the one ordered list lives in a guide section. A blockquote
// does occur in sections, but its content is the same markup as a paragraph and
// adds no meaning, so the decoder flattens it into paragraphs rather than
// modeling it.
//
// Inside an entity section inline emphasis never nests — no bold wraps an
// italic, no link wraps more than a single run — so inline content is flat: a
// [Text] or [Link] carries one [Style] rather than a subtree of inline nodes.
package prose

// Block represents a block-level node in a prose tree. The concrete variants
// are [Paragraph] and [List].
type Block interface {
	isBlock()
}

// Inline represents an inline-level node within a block or phrase. The concrete
// variants are [Text], [Link], and [LineBreak].
type Inline interface {
	isInline()
}

// Passage represents prose as a sequence of blocks.
type Passage struct {
	blocks []Block
}

// NewPassage constructs a passage from a sequence of blocks.
func NewPassage(blocks ...Block) Passage {
	return Passage{blocks: blocks}
}

// Blocks returns the blocks of the passage.
func (p Passage) Blocks() []Block {
	return p.blocks
}

// Phrase represents prose as a sequence of inline runs, with no block
// structure. It is the inline-only counterpart of a [Passage], carrying the
// text of a table cell — a field type or description — that holds no paragraphs
// or lists.
type Phrase struct {
	inlines []Inline
}

// NewPhrase constructs a phrase from inline content.
func NewPhrase(inlines ...Inline) Phrase {
	return Phrase{inlines: inlines}
}

// Inlines returns the inline content of the phrase.
func (p Phrase) Inlines() []Inline {
	return p.inlines
}
