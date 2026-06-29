// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pipeline

import (
	"fmt"
)

// Operator produces a table from one or more input tables. Constructing an
// operator only wires its inputs; Apply runs the operation, returning the
// resulting table or the failure that prevented it.
type Operator[K comparable, R Record] interface {
	Apply() (Table[K, R], error)
}

// Mapping transforms one record into another. It fails when the record cannot
// be transformed.
type Mapping[A, B Record] interface {
	Apply(record A) (B, error)
}

// MappedTable is the projection operator: it applies a mapping to every record
// of a source table, carrying each record's key through unchanged.
type MappedTable[K comparable, A, B Record] struct {
	source  Table[K, A]
	mapping Mapping[A, B]
}

// NewMappedTable constructs the projection of source through mapping.
func NewMappedTable[K comparable, A, B Record](
	source Table[K, A], mapping Mapping[A, B],
) MappedTable[K, A, B] {
	return MappedTable[K, A, B]{source: source, mapping: mapping}
}

// Apply returns the projected table, one record per source record under the
// same key. It fails when the mapping fails on any record.
func (t MappedTable[K, A, B]) Apply() (Table[K, B], error) {
	out := NewMapTableWithCapacity[K, B](t.source.Count())
	for key, record := range t.source.All() {
		mapped, err := t.mapping.Apply(record)
		if err != nil {
			return out, fmt.Errorf("mapping record %v: %w", key, err)
		}
		out.Insert(key, mapped)
	}
	return out, nil
}

// MergedTable is the union operator: it gathers every record from two tables
// whose key sets are disjoint.
type MergedTable[K comparable, R Record] struct {
	left  Table[K, R]
	right Table[K, R]
}

// NewMergedTable constructs the union of left and right.
func NewMergedTable[K comparable, R Record](left, right Table[K, R]) MergedTable[K, R] {
	return MergedTable[K, R]{left: left, right: right}
}

// Apply returns the union of the two tables. It fails when a key appears in
// both, since a key identifies at most one record.
func (t MergedTable[K, R]) Apply() (Table[K, R], error) {
	out := NewMapTableWithCapacity[K, R](t.left.Count() + t.right.Count())
	for _, source := range []Table[K, R]{t.left, t.right} {
		for key, record := range source.All() {
			if _, exists := out.Lookup(key); exists {
				return out, fmt.Errorf("merging tables: key %v present in both", key)
			}
			out.Insert(key, record)
		}
	}
	return out, nil
}
