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

// All decodes the nodes into inline prose content. Emphasis collapses into a
// styled run, a break into a line break, an emoji image into its alt text. It
// returns an error for markup the flat model cannot hold.
func (i Inlines) All() ([]prose.Inline, error) {
	var out []prose.Inline
	for node := range i.nodes.All() {
		switch node.Tag() {
		case TagText:
			if raw := node.Text(); raw != "" {
				out = append(out, prose.NewText(raw, prose.StylePlain))
			}
		case TagEm:
			run, err := i.styled(node, prose.StyleItalic)
			if err != nil {
				return nil, err
			}
			out = append(out, run)
		case TagStrong:
			run, err := i.styled(node, prose.StyleBold)
			if err != nil {
				return nil, err
			}
			out = append(out, run)
		case TagCode:
			run, err := i.styled(node, prose.StyleCode)
			if err != nil {
				return nil, err
			}
			out = append(out, run)
		case TagA:
			run, err := i.link(node)
			if err != nil {
				return nil, err
			}
			out = append(out, run...)
		case TagBr:
			out = append(out, prose.NewLineBreak())
		case TagImg:
			run, err := i.image(node)
			if err != nil {
				return nil, err
			}
			out = append(out, run)
		default:
			return nil, fmt.Errorf("parsing inline: unexpected markup %q", node.Text())
		}
	}
	return out, nil
}

// styled decodes a leaf emphasis element into a single styled run. It returns an
// error when the element nests further markup, which the flat model cannot hold.
func (i Inlines) styled(node Node, style prose.Style) (prose.Inline, error) {
	for range node.Children().All() {
		return nil, fmt.Errorf("parsing inline: nested markup in styled run %q", node.Text())
	}
	return prose.NewText(node.Text(), style), nil
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
func (i Inlines) image(node Node) (prose.Inline, error) {
	if !node.HasClass("emoji") {
		return nil, fmt.Errorf("parsing inline: non-emoji image %q", node.Attr("src"))
	}
	return prose.NewText(node.Attr("alt"), prose.StylePlain), nil
}
