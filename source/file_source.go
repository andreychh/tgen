// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package source

import (
	"context"
	"fmt"
	"io"
	"os"
)

// FileSource provides access to local filesystem resources.
type FileSource struct {
	name string
}

// NewFileSource creates a FileSource for the specified file path.
func NewFileSource(name string) FileSource {
	return FileSource{name: name}
}

// Open prepares the file for reading. The context is ignored.
func (s FileSource) Open(_ context.Context) (io.ReadCloser, error) {
	file, err := os.Open(s.name)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", s.name, err)
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("stating file %q: %w", s.name, err)
	}
	if info.IsDir() {
		_ = file.Close()
		return nil, fmt.Errorf("expected file, but %q is a directory", s.name)
	}
	return file, nil
}
