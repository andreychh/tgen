// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

import "time"

// Snapshot captures binary metadata and the moment of code generation.
//
// It is designed to be passed to templates for generating file headers. Use
// [NewSnapshot] to create one at the start of a generation run, then call
// Elapsed at the end to measure the duration.
type Snapshot struct {
	meta Meta
	at   time.Time
}

// NewSnapshot creates a Snapshot using the provided Meta and the current time.
func NewSnapshot(m Meta) Snapshot {
	return NewSnapshotAt(m, time.Now())
}

// NewSnapshotAt creates a Snapshot with an explicitly provided generation time.
func NewSnapshotAt(m Meta, t time.Time) Snapshot {
	return Snapshot{meta: m, at: t}
}

// Meta returns the binary metadata associated with this snapshot.
func (s Snapshot) Meta() Meta {
	return s.meta
}

// CreatedAt returns the time at which the snapshot was created.
func (s Snapshot) CreatedAt() time.Time {
	return s.at
}

// Elapsed returns the duration since the snapshot was created.
func (s Snapshot) Elapsed() time.Duration {
	return time.Since(s.at)
}
