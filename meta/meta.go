// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package meta provides access to binary metadata embedded at build time and
// read from the runtime environment.
//
// Metadata is organized into three domains: [VCS] for version control
// information, [Release] for versioning and build origin, and [Build] for
// compiler and platform details. All three are accessible through [Meta].
//
// The source of metadata is abstracted behind the [Source] interface. Use
// [NewDetectedSource] to automatically infer the build scenario at runtime, or
// supply a specific implementation for testing.
package meta

// Source provides access to binary metadata by key.
//
// Implementations are responsible for mapping well-known keys (see
// source_keys.go) to their respective data sources.
type Source interface {
	// Get returns the value and true if the key is known, or an empty string and
	// false if it is not.
	Get(key string) (value string, exists bool)
}

// Meta provides metadata about the binary.
type Meta struct {
	source Source
}

// NewMeta creates a Meta backed by the provided Source.
func NewMeta(source Source) Meta {
	return Meta{source: source}
}

// VCS returns version control metadata.
func (m Meta) VCS() VCS { return NewVCS(m.source) }

// Release returns release metadata.
func (m Meta) Release() Release { return NewRelease(m.source) }

// Build returns compiler and platform metadata.
func (m Meta) Build() Build { return NewBuild(m.source) }
