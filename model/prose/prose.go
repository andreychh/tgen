// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package prose models Telegram documentation text as a target-neutral tree of
// block and inline nodes, decoupled from both the source HTML and any output
// target. The decoder maps HTML into this tree; each target renders the tree
// into its own doc-comment syntax. The same tree carries any free-form cell —
// a field description as much as a field type — so the model names the text,
// not its role.
//
// The tree models only the prose found inside entity sections — objects,
// methods, and unions. Markup that appears elsewhere on the page (the intro,
// the changelog, standalone guide sections) is out of scope, so constructs
// that occur only there have no node: there is no image, because real images
// live in the intro, and no ordered list, because the one ordered list lives in
// a guide section.
package prose

// Block represents a block-level node in a prose tree. The concrete variants
// are Paragraph, CodeBlock, and List.
type Block interface {
	isBlock()
}

// Inline represents an inline-level node nested within a block. The concrete
// variants are Text, Bold, Italic, Code, Link, and LineBreak.
type Inline interface {
	isInline()
}

// Tree represents a prose tree as a sequence of blocks.
type Tree struct {
	blocks []Block
}

// NewTree constructs a prose tree from a sequence of blocks.
func NewTree(blocks ...Block) Tree {
	return Tree{blocks: blocks}
}

// Blocks returns the blocks of the prose tree.
func (t Tree) Blocks() []Block {
	return t.blocks
}
