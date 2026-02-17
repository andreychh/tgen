// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering

import (
	"fmt"
	"os"
	"path/filepath"
)

type Fileset struct {
	artifacts Artifacts
}

func NewFileset(a Artifacts) Fileset {
	return Fileset{artifacts: a}
}

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
