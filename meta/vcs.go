// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

import "time"

// TreeState indicates the state of the working tree at build time.
type TreeState string

const (
	// TreeStateUnknown is returned when the tree state is not available.
	TreeStateUnknown TreeState = "unknown"
	// TreeStateClean indicates no uncommitted changes were present.
	TreeStateClean TreeState = "clean"
	// TreeStateDirty indicates uncommitted changes were present.
	TreeStateDirty TreeState = "dirty"
)

// VCS provides version control metadata for the commit the binary was built
// from.
type VCS struct {
	source Source
}

// NewVCS creates a VCS backed by the provided Source.
func NewVCS(source Source) VCS {
	return VCS{source: source}
}

// Revision returns the commit identifier.
func (v VCS) Revision() Revision {
	return NewRevision(v.source)
}

// Date returns the commit timestamp, or the zero value if unavailable.
func (v VCS) Date() time.Time {
	value, ok := v.source.Get(KeyVCSTime)
	if !ok {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return t
}

// TreeState returns the state of the working tree at build time.
func (v VCS) TreeState() TreeState {
	value, ok := v.source.Get(KeyVCSTreeState)
	if !ok {
		return TreeStateUnknown
	}
	if value == "true" {
		return TreeStateDirty
	}
	return TreeStateClean
}
