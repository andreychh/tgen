// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

import (
	"iter"

	"github.com/PuerkitoBio/goquery"
)

// Node represents a single source node — an element or a text node. It answers
// structural questions about itself and knows nothing of prose.
type Node struct {
	sel *goquery.Selection
}

// NewNode constructs a Node over a single source selection.
func NewNode(sel *goquery.Selection) Node {
	return Node{sel: sel}
}

// Tag returns the node's tag, or [TagUnknown] for a tag outside entity-section
// prose. A text node is [TagText].
func (n Node) Tag() Tag {
	switch goquery.NodeName(n.sel) {
	case "#text":
		return TagText
	case "em":
		return TagEm
	case "strong":
		return TagStrong
	case "code":
		return TagCode
	case "a":
		return TagA
	case "br":
		return TagBr
	case "img":
		return TagImg
	case "p":
		return TagP
	case "ul":
		return TagUl
	case "li":
		return TagLi
	case "blockquote":
		return TagBlockquote
	case "div":
		return TagDiv
	default:
		return TagUnknown
	}
}

// Children returns the node's element children, excluding text nodes.
func (n Node) Children() Nodes {
	return NewNodes(n.sel.Children())
}

// Contents returns the node's children, text and elements alike.
func (n Node) Contents() Nodes {
	return NewNodes(n.sel.Contents())
}

// Text returns the verbatim text content of the node.
func (n Node) Text() string {
	return n.sel.Text()
}

// Attr returns the value of the named attribute, or the empty string.
func (n Node) Attr(name string) string {
	value, _ := n.sel.Attr(name)
	return value
}

// HasClass reports whether the node carries the given class.
func (n Node) HasClass(class string) bool {
	return n.sel.HasClass(class)
}

// Nodes represents a sequence of sibling source nodes.
type Nodes struct {
	sel *goquery.Selection
}

// NewNodes constructs a Nodes over a multi-node selection.
func NewNodes(sel *goquery.Selection) Nodes {
	return Nodes{sel: sel}
}

// All iterates the sequence in document order.
func (ns Nodes) All() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, inner := range ns.sel.EachIter() {
			if !yield(NewNode(inner)) {
				return
			}
		}
	}
}
