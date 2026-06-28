// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

// Cursor is a one-pass reader over a slice of items, consuming them one at a
// time from the front.
type Cursor[T any] struct {
	items []T
	pos   int
}

// NewCursor constructs a Cursor positioned before the first item.
func NewCursor[T any](items []T) *Cursor[T] {
	return &Cursor[T]{
		items: items,
		pos:   0,
	}
}

// Peek returns the current item without consuming it. The boolean reports
// whether an item was available; on an exhausted cursor it is false and the
// item is the zero value.
func (c *Cursor[T]) Peek() (T, bool) {
	if c.Done() {
		var zero T
		return zero, false
	}
	return c.items[c.pos], true
}

// Take consumes the current item and returns it. The boolean reports whether an
// item was available; on an exhausted cursor it is false and the item is the
// zero value.
func (c *Cursor[T]) Take() (T, bool) {
	if c.Done() {
		var zero T
		return zero, false
	}
	item := c.items[c.pos]
	c.pos++
	return item, true
}

// Skip consumes the current item without returning it, advancing the cursor.
func (c *Cursor[T]) Skip() {
	c.pos++
}

// Done reports whether the cursor has no items left to consume.
func (c *Cursor[T]) Done() bool {
	return c.pos >= len(c.items)
}
