// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package buildinfo

import "time"

// Snapshot captures the tool's build metadata at a specific point in time.
//
// It is designed to be used in templates for generating file headers (e.g.,
// copyright notices or generation timestamps).
type Snapshot struct {
	info *Info
	date time.Time
}

// NewSnapshotAt creates a Snapshot with explicitly provided build information
// and time.
func NewSnapshotAt(info *Info, date time.Time) Snapshot {
	return Snapshot{
		info: info,
		date: date,
	}
}

// NewSnapshot creates a new Snapshot using the provided build information and
// the current local time.
func NewSnapshot(info *Info) Snapshot {
	return Snapshot{
		info: info,
		date: time.Now(),
	}
}

// Version returns the semantic version of the tool (e.g., "v1.0.0").
func (s Snapshot) Version() string {
	return s.info.Version
}

// GenerationDate returns the snapshot's date formatted as "YYYY-MM-DD".
func (s Snapshot) GenerationDate() string {
	return s.date.Format(time.DateOnly)
}
