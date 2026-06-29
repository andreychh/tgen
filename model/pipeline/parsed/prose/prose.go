// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package prose decodes the documentation source of an entity section into the
// flat prose model: a passage of blocks for a section body, a phrase of runs
// for a table cell. It maps the closed set of HTML nodes that occur in such
// sections, and rejects any markup the model cannot represent.
package prose

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model/prose"
)

// Passage represents the block-level prose of a documentation section body,
// carried by a sequence of source nodes.
type Passage struct {
	nodes Nodes
}

// NewPassage constructs a Passage over a sequence of sibling source nodes.
func NewPassage(sel *goquery.Selection) Passage {
	return Passage{nodes: NewNodes(sel)}
}

// Value decodes the nodes into a prose passage. It returns an error when the
// source carries markup the flat model cannot hold: nested inline emphasis, a
// multi-run link, a non-emoji image, or an unexpected tag.
func (p Passage) Value() (prose.Passage, error) {
	blocks, err := NewBlocks(p.nodes).All()
	if err != nil {
		return prose.Passage{}, fmt.Errorf("parsing prose: %w", err)
	}
	return prose.NewPassage(blocks...), nil
}

// Phrase represents the inline-level prose of a table cell, carried by a
// sequence of source nodes.
type Phrase struct {
	nodes Nodes
}

// NewPhrase constructs a Phrase over a sequence of sibling source nodes.
func NewPhrase(sel *goquery.Selection) Phrase {
	return Phrase{nodes: NewNodes(sel)}
}

// Value decodes the nodes into a prose phrase. It returns an error when the
// source carries markup the flat model cannot hold: nested inline emphasis, a
// multi-run link, a non-emoji image, or unexpected markup.
func (p Phrase) Value() (prose.Phrase, error) {
	inlines, err := NewInlines(p.nodes).All()
	if err != nil {
		return prose.Phrase{}, fmt.Errorf("parsing prose: %w", err)
	}
	return prose.NewPhrase(inlines...), nil
}
