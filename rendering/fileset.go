// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering

import (
	"fmt"
	"os"
	"path/filepath"
)

// Fileset represents a collection of [Artifacts] ready to be materialized onto
// the file system.
type Fileset struct {
	artifacts Artifacts
}

// NewFileset returns a [Fileset] containing the provided [Artifacts].
func NewFileset(a Artifacts) Fileset {
	return Fileset{artifacts: a}
}

// Emit writes all artifacts in the [Fileset] to the specified output directory.
// It creates the directory and any necessary parents if they do not exist.
func (f Fileset) Emit(path string) error {
	err := os.MkdirAll(path, 0o750)
	if err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	for filename, view := range f.artifacts {
		err := f.renderFile(filepath.Join(path, filename), view)
		if err != nil {
			return fmt.Errorf("rendering file %q: %w", filename, err)
		}
	}
	return nil
}

// renderFile creates and safely opens a target file, delegating the actual
// content generation to the provided [View].
func (f Fileset) renderFile(path string, v View) error {
	file, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("creating file %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	err = v.Render(file)
	if err != nil {
		return fmt.Errorf("rendering view to %q: %w", path, err)
	}
	return nil
}
