// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

import (
	"fmt"

	"github.com/andreychh/tgen/model/prose"
)

// Inlines represents the inline content carried by a sequence of source nodes,
// decoded on demand into flat inline runs.
type Inlines struct {
	nodes Nodes
}

// NewInlines constructs an Inlines over a sequence of sibling source nodes.
func NewInlines(nodes Nodes) Inlines {
	return Inlines{nodes: nodes}
}

// All decodes the nodes into inline prose content. It returns an error for
// markup the flat model cannot hold.
func (i Inlines) All() ([]prose.Inline, error) {
	var out []prose.Inline
	for node := range i.nodes.All() {
		runs, err := i.runs(node)
		if err != nil {
			return nil, err
		}
		out = append(out, runs...)
	}
	return out, nil
}

// runs decodes a single node into the inline runs it contributes. Emphasis
// collapses into a styled run, a break into a line break, an emoji image into
// its alt text. It returns an error for markup the flat model cannot hold.
func (i Inlines) runs(node Node) ([]prose.Inline, error) {
	switch node.Tag() {
	case TagText:
		return i.text(node), nil
	case TagEm:
		return i.styled(node, prose.StyleItalic)
	case TagStrong:
		return i.styled(node, prose.StyleBold)
	case TagCode:
		return i.styled(node, prose.StyleCode)
	case TagA:
		return i.link(node)
	case TagBr:
		return []prose.Inline{prose.NewLineBreak()}, nil
	case TagImg:
		return i.image(node)
	case TagUnknown, TagP, TagUl, TagLi, TagBlockquote, TagDiv:
		return nil, fmt.Errorf("parsing inline: unexpected markup %q", node.Text())
	}
	return nil, nil
}

// text decodes a text node into its run, or nothing when the text is empty.
func (i Inlines) text(node Node) []prose.Inline {
	if raw := node.Text(); raw != "" {
		return []prose.Inline{prose.NewText(raw, prose.StylePlain)}
	}
	return nil
}

// styled decodes a leaf emphasis element into a single styled run. It returns
// an error when the element nests further markup, which the flat model cannot
// hold.
func (i Inlines) styled(node Node, style prose.Style) ([]prose.Inline, error) {
	for range node.Children().All() {
		return nil, fmt.Errorf("parsing inline: nested markup in styled run %q", node.Text())
	}
	return []prose.Inline{prose.NewText(node.Text(), style)}, nil
}

// link decodes an anchor into a single styled run carrying its href. An empty
// anchor yields nothing; an anchor holding more than one run or a non-text run
// returns an error.
func (i Inlines) link(node Node) ([]prose.Inline, error) {
	runs, err := NewInlines(node.Contents()).All()
	if err != nil {
		return nil, err
	}
	if len(runs) == 0 {
		return nil, nil
	}
	if len(runs) != 1 {
		return nil, fmt.Errorf("parsing inline: link is not a single run %q", node.Text())
	}
	run, ok := runs[0].(prose.Text)
	if !ok {
		return nil, fmt.Errorf("parsing inline: link wraps non-text %q", node.Text())
	}
	return []prose.Inline{prose.NewLink(run.Content(), run.Style(), node.Attr("href"))}, nil
}

// image decodes an emoji image into its alt text. It returns an error for any
// non-emoji image, whose place is the intro rather than entity-section prose.
func (i Inlines) image(node Node) ([]prose.Inline, error) {
	if !node.HasClass("emoji") {
		return nil, fmt.Errorf("parsing inline: non-emoji image %q", node.Attr("src"))
	}
	return []prose.Inline{prose.NewText(node.Attr("alt"), prose.StylePlain)}, nil
}
