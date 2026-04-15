// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package output

import (
	"fmt"
	"os"
	"path/filepath"
)

// Fileset represents a set of [Artifacts] that can be written to the file
// system.
type Fileset struct {
	artifacts Artifacts
}

// NewFileset constructs a [Fileset] from the provided [Artifacts].
func NewFileset(a Artifacts) Fileset {
	return Fileset{artifacts: a}
}

// Emit writes all artifacts to path, creating it and any missing parents if
// they do not exist. Returns an error if the directory cannot be created or any
// artifact fails to render.
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

// renderFile writes view to path, creating the file and any missing parent
// directories if needed.
func (f Fileset) renderFile(path string, view View) error {
	clean := filepath.Clean(path)
	err := os.MkdirAll(filepath.Dir(clean), 0o750)
	if err != nil {
		return fmt.Errorf("creating parent directory: %w", err)
	}
	file, err := os.Create(clean)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	err = view.Render(file)
	if err != nil {
		return fmt.Errorf("rendering view to %q: %w", path, err)
	}
	return nil
}
