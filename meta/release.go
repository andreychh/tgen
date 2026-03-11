// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

// Release provides metadata about the versioned release of the binary.
type Release struct {
	source Source
}

// NewRelease creates a Release backed by the provided Source.
func NewRelease(source Source) Release {
	return Release{source: source}
}

// Version returns the semantic version of the release, or "unknown" if
// unavailable.
func (r Release) Version() string {
	value, ok := r.source.Get(KeyVersion)
	if !ok {
		return unknownValue
	}
	return value
}

// Builder returns the name of the tool that produced the binary, or "unknown" if
// unavailable.
func (r Release) Builder() string {
	value, ok := r.source.Get(KeyBuilder)
	if !ok {
		return unknownValue
	}
	return value
}
