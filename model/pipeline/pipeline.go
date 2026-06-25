// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pipeline

import (
	"iter"
)

// Record represents a single row held by a [Table].
type Record any

// Table represents a set of records keyed by K, each key identifying at most
// one record.
type Table[K comparable, R Record] interface {
	// Insert adds record under key. The key must be absent; inserting a key that is
	// already present panics.
	Insert(key K, record R)
	// Lookup returns the record stored under key. The boolean reports whether such
	// a record is present; when absent, the record is the zero value.
	Lookup(key K) (R, bool)
	// Count returns the number of records in the table.
	Count() int
	// All returns an iterator over every key and its record, in unspecified order.
	All() iter.Seq2[K, R]
}
