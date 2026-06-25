// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

import (
	"fmt"

	"github.com/andreychh/tgen/model/prose"
)

// Blocks represents the block content carried by a sequence of source nodes,
// decoded on demand into the flat prose model.
type Blocks struct {
	nodes Nodes
}

// NewBlocks constructs a Blocks over a sequence of sibling source nodes.
func NewBlocks(nodes Nodes) Blocks {
	return Blocks{nodes: nodes}
}

// All decodes the nodes into block prose content. It returns an error for
// markup the flat model cannot hold.
func (b Blocks) All() ([]prose.Block, error) {
	var out []prose.Block
	for node := range b.nodes.All() {
		blocks, err := b.blocks(node)
		if err != nil {
			return nil, err
		}
		out = append(out, blocks...)
	}
	return out, nil
}

// blocks decodes a single node into the blocks it contributes. A paragraph and
// a list contribute one block each; a blockquote contributes its own blocks; a
// blog-image div contributes nothing. It returns an error for any tag the flat
// model cannot hold.
func (b Blocks) blocks(node Node) ([]prose.Block, error) {
	switch node.Tag() {
	case TagP:
		return b.paragraph(node)
	case TagUl:
		return b.list(node)
	case TagBlockquote:
		return NewBlocks(node.Children()).All()
	case TagDiv:
		return b.division(node)
	case TagText, TagEm, TagStrong, TagCode, TagA, TagBr, TagImg, TagLi, TagUnknown:
		return nil, fmt.Errorf("parsing block: unexpected tag %q", node.Text())
	}
	return nil, nil
}

// paragraph decodes a <p> node into a single paragraph block. It returns an
// error for inline markup the flat model cannot hold.
func (b Blocks) paragraph(node Node) ([]prose.Block, error) {
	inlines, err := NewInlines(node.Contents()).All()
	if err != nil {
		return nil, err
	}
	return []prose.Block{prose.NewParagraph(inlines...)}, nil
}

// list decodes a <ul> node into a single list block. It returns an error for
// any child that is not a list item.
func (b Blocks) list(node Node) ([]prose.Block, error) {
	var items []prose.Item
	for child := range node.Children().All() {
		if child.Tag() != TagLi {
			return nil, fmt.Errorf("parsing list: unexpected tag %q", child.Text())
		}
		inlines, err := NewInlines(child.Contents()).All()
		if err != nil {
			return nil, err
		}
		items = append(items, prose.NewItem(inlines...))
	}
	return []prose.Block{prose.NewList(items...)}, nil
}

// division decodes a <div> node, which contributes blocks only as a blog-image
// wrapper carrying none. It returns an error for any other div.
func (b Blocks) division(node Node) ([]prose.Block, error) {
	if !node.HasClass("blog_image_wrap") {
		return nil, fmt.Errorf("parsing block: unexpected div %q", node.Text())
	}
	return nil, nil
}
