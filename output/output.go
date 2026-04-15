// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package output turns structured views into source files on disk.
package output

import (
	"io"
)

// View represents a unit of generated output that can be written to an
// [io.Writer].
type View interface {
	// Render writes the generated content to w. Returns an error if writing fails.
	Render(w io.Writer) error
}

// Artifacts maps output file names to their corresponding [View].
type Artifacts map[string]View
