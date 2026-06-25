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

// All decodes the nodes into block prose content. A blockquote contributes its
// own blocks; a blog-image div contributes nothing. It returns an error for any
// tag the flat model cannot hold.
func (b Blocks) All() ([]prose.Block, error) {
	var out []prose.Block
	for node := range b.nodes.All() {
		switch node.Tag() {
		case TagP:
			inlines, err := NewInlines(node.Contents()).All()
			if err != nil {
				return nil, err
			}
			out = append(out, prose.NewParagraph(inlines...))
		case TagUl:
			list, err := NewBlocks(node.Children()).list()
			if err != nil {
				return nil, err
			}
			out = append(out, list)
		case TagBlockquote:
			inner, err := NewBlocks(node.Children()).All()
			if err != nil {
				return nil, err
			}
			out = append(out, inner...)
		case TagDiv:
			if !node.HasClass("blog_image_wrap") {
				return nil, fmt.Errorf("parsing block: unexpected div %q", node.Text())
			}
		default:
			return nil, fmt.Errorf("parsing block: unexpected tag %q", node.Text())
		}
	}
	return out, nil
}

// list decodes the items of a <ul> into a prose list. It returns an error for
// any child that is not a list item.
func (b Blocks) list() (prose.List, error) {
	var items []prose.Item
	for node := range b.nodes.All() {
		if node.Tag() != TagLi {
			return prose.List{}, fmt.Errorf("parsing list: unexpected tag %q", node.Text())
		}
		inlines, err := NewInlines(node.Contents()).All()
		if err != nil {
			return prose.List{}, err
		}
		items = append(items, prose.NewItem(inlines...))
	}
	return prose.NewList(items...), nil
}
