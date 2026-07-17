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

// PartialMapping transforms one record into another, or reports that no
// record results. It fails when the record cannot be evaluated.
type PartialMapping[A, B Record] interface {
	Apply(record A) (B, bool, error)
}

// PartialMappedTable is the filter-projection operator: it applies a partial
// mapping to every record of a source table, keeping only the records for
// which the mapping produces one, and carrying each kept record's key through
// unchanged.
type PartialMappedTable[K comparable, A, B Record] struct {
	source  Table[K, A]
	mapping PartialMapping[A, B]
}

// NewPartialMappedTable constructs the filtered projection of source through
// mapping.
func NewPartialMappedTable[K comparable, A, B Record](
	source Table[K, A], mapping PartialMapping[A, B],
) PartialMappedTable[K, A, B] {
	return PartialMappedTable[K, A, B]{source: source, mapping: mapping}
}

// Apply returns the projected table, one record per source record for which
// mapping produced one, under the same key. It fails when the mapping fails on
// any record.
func (t PartialMappedTable[K, A, B]) Apply() (Table[K, B], error) {
	out := NewMapTableWithCapacity[K, B](t.source.Count())
	for key, record := range t.source.All() {
		mapped, ok, err := t.mapping.Apply(record)
		if err != nil {
			return out, fmt.Errorf("mapping record %v: %w", key, err)
		}
		if !ok {
			continue
		}
		out.Insert(key, mapped)
	}
	return out, nil
}

// Filter decides whether a record, held under key, belongs in a filtered
// table. It never fails: every key and record combination yields a definite
// yes or no.
type Filter[K comparable, R Record] interface {
	Apply(key K, record R) bool
}

// FilteredTable is the restriction operator: it keeps every record of source
// for which filter reports true, discarding the rest under their original
// keys.
type FilteredTable[K comparable, R Record] struct {
	source Table[K, R]
	filter Filter[K, R]
}

// NewFilteredTable constructs the restriction of source to the records filter
// keeps.
func NewFilteredTable[K comparable, R Record](
	source Table[K, R], filter Filter[K, R],
) FilteredTable[K, R] {
	return FilteredTable[K, R]{source: source, filter: filter}
}

// Apply returns the restriction of source to the records filter keeps, one
// per kept record under its original key.
func (t FilteredTable[K, R]) Apply() Table[K, R] {
	out := NewMapTableWithCapacity[K, R](t.source.Count())
	for key, record := range t.source.All() {
		if !t.filter.Apply(key, record) {
			continue
		}
		out.Insert(key, record)
	}
	return out
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
