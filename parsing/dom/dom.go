// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package dom provides HTML traversal and data extraction primitives.
package dom

import "iter"

// Selection represents a set of HTML nodes.
//
// For all text extraction methods (e.g. Text, Attr) whitespace normalization
// MUST be performed: trimming leading/trailing spaces and collapsing internal
// whitespace sequences to a single space.
//
// Methods return a new Selection, leaving the original unmodified.
type Selection interface {
	// Text returns the combined text content of the nodes and their descendants.
	//
	// Note: The text is concatenated directly without any separators (spaces or
	// newlines) between elements. For example, "<li>A</li><li>B</li>" results in
	// "AB".
	Text() string

	// Attr returns the attribute value of the first node in the set. The boolean
	// reports whether the attribute exists.
	Attr(name string) (value string, exists bool)

	// First returns a new Selection containing only the first node.
	First() Selection

	// Find returns a new Selection containing descendants matching the selector.
	Find(selector string) Selection

	// Filter reduces the set to nodes matching the selector.
	Filter(selector string) Selection

	// FilterFunc reduces the set to nodes satisfying the predicate f.
	FilterFunc(f func(Selection) bool) Selection

	// NextUntil returns following siblings up to (but not including) the selector.
	NextUntil(selector string) Selection

	// Length returns the number of nodes in the set.
	Length() int

	// IsEmpty reports whether the set contains no nodes.
	IsEmpty() bool

	// At returns a new Selection containing the node at the specified index. It
	// returns an empty selection if the index is negative or greater than the
	// number of elements.
	At(index int) Selection

	// All returns an iterator over the nodes in the set.
	All() iter.Seq2[int, Selection]
}
