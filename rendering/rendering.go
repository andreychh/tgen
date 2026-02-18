// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package rendering provides an abstraction layer for code generation.
package rendering

import (
	"io"
)

// View represents a generic renderable component that writes generated content
// to an [io.Writer], abstracting away the underlying templating mechanism.
type View interface {
	// Render writes the generated output to w.
	Render(w io.Writer) error
}

// Artifacts maps target file paths to their corresponding [View]. It acts as a
// virtual, in-memory representation of the files to be generated.
type Artifacts map[string]View
