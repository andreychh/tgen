// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

// Build provides compiler and platform metadata.
type Build struct {
	source Source
}

// NewBuild creates a Build backed by the provided Source.
func NewBuild(source Source) Build {
	return Build{source: source}
}

// GoVersion returns the version of the Go toolchain used to build the binary,
// or "unknown" if unavailable.
func (b Build) GoVersion() string {
	value, ok := b.source.Get(KeyGoVersion)
	if !ok {
		return "unknown"
	}
	return value
}

// Platform returns the target operating system and architecture (e.g.
// "linux/amd64"), or "unknown" if unavailable.
func (b Build) Platform() string {
	value, ok := b.source.Get(KeyPlatform)
	if !ok {
		return "unknown"
	}
	return value
}
