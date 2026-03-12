// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering

import "github.com/andreychh/tgen/meta"

// Snapshot represents the generation moment metadata, adapted for use in
// templates.
type Snapshot struct {
	inner meta.Snapshot
}

// NewSnapshot creates a Snapshot from the given [meta.Snapshot].
func NewSnapshot(s meta.Snapshot) Snapshot {
	return Snapshot{inner: s}
}

// Date returns the generation date formatted as "YYYY-MM-DD".
func (s Snapshot) Date() string {
	return s.inner.CreatedAt().Format("2006-01-02")
}

// Meta returns the binary metadata.
func (s Snapshot) Meta() Meta {
	return NewMeta(s.inner.Meta())
}
