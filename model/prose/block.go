// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

// Paragraph represents a run of inline content delimited as one block.
type Paragraph struct {
	inlines []Inline
}

// NewParagraph constructs a paragraph from inline content.
func NewParagraph(inlines ...Inline) Paragraph {
	return Paragraph{inlines: inlines}
}

// Inlines returns the inline content of the paragraph.
func (p Paragraph) Inlines() []Inline {
	return p.inlines
}

func (Paragraph) isBlock() {}

// List represents a sequence of items. Entity-section prose carries only
// unordered lists, so the tree models no ordering.
type List struct {
	items []Item
}

// NewList constructs a list from its items.
func NewList(items ...Item) List {
	return List{items: items}
}

// Items returns the items of the list.
func (l List) Items() []Item {
	return l.items
}

func (List) isBlock() {}

// Item represents a single list entry of inline content.
type Item struct {
	inlines []Inline
}

// NewItem constructs a list item from inline content.
func NewItem(inlines ...Inline) Item {
	return Item{inlines: inlines}
}

// Inlines returns the inline content of the item.
func (i Item) Inlines() []Inline {
	return i.inlines
}
