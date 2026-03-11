// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

// Revision is a VCS commit identifier.
type Revision struct {
	source Source
}

// NewRevision creates a Revision backed by the provided Source.
func NewRevision(source Source) Revision {
	return Revision{source: source}
}

// Full returns the full commit hash, or "unknown" if unavailable.
func (r Revision) Full() string {
	value, ok := r.source.Get(KeyVCSRevision)
	if !ok {
		return unknownValue
	}
	return value
}

// Short returns the first 7 characters of the commit hash, or "unknown" if
// unavailable. It panics if the hash is present but shorter than 7 characters.
func (r Revision) Short() string {
	value, ok := r.source.Get(KeyVCSRevision)
	if !ok {
		return unknownValue
	}
	return value[:7]
}
