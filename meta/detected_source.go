// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

// DetectedSource infers the build scenario at runtime and delegates to the
// corresponding [Source] implementation.
type DetectedSource struct{}

// NewDetectedSource creates a DetectedSource.
func NewDetectedSource() DetectedSource {
	return DetectedSource{}
}

// Get returns the metadata value for the given key.
func (DetectedSource) Get(key string) (string, bool) {
	if version != unknownValue {
		return NewLDFlagsSource().Get(key)
	}
	return NewRuntimeSource().Get(key)
}
