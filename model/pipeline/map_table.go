// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pipeline

import (
	"fmt"
	"iter"
	"maps"
)

// MapTable is the in-memory implementation of [Table]. Its zero value is not
// usable; construct one with [NewMapTable] or [NewMapTableWithCapacity]. A
// MapTable is not safe for concurrent use.
type MapTable[K comparable, R Record] struct {
	inner map[K]R
}

// NewMapTableWithCapacity constructs an empty MapTable with room preallocated
// for cap records.
func NewMapTableWithCapacity[K comparable, R Record](cap int) MapTable[K, R] {
	return MapTable[K, R]{
		inner: make(map[K]R, cap),
	}
}

// NewMapTable constructs an empty MapTable.
func NewMapTable[K comparable, R Record]() MapTable[K, R] {
	return NewMapTableWithCapacity[K, R](0)
}

func (m MapTable[K, R]) Insert(key K, record R) {
	if _, exists := m.inner[key]; exists {
		panic(fmt.Sprintf("MapTable.Insert: key %v already present", key))
	}
	m.inner[key] = record
}

func (m MapTable[K, R]) Lookup(key K) (R, bool) {
	row, exists := m.inner[key]
	return row, exists
}

func (m MapTable[K, R]) Count() int {
	return len(m.inner)
}

func (m MapTable[K, R]) All() iter.Seq2[K, R] {
	return maps.All(m.inner)
}
