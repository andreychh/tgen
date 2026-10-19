// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package output

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
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

// Emit writes every artifact to path once all of them have rendered, creating
// path and any missing parents. If an artifact fails to render, Emit returns an
// error and leaves the file system untouched. It also returns an error if a
// directory or a file cannot be written; files written before that failure stay
// in place.
func (f Fileset) Emit(path string) error {
	batch, err := f.render()
	if err != nil {
		return err
	}
	return batch.write(path)
}

// render returns every artifact rendered, without touching the file system. It
// returns an error for the first artifact in name order whose view fails, so
// the error is the same on every run.
func (f Fileset) render() (batch, error) {
	result := make(batch, len(f.artifacts))
	for _, name := range slices.Sorted(maps.Keys(f.artifacts)) {
		var buf bytes.Buffer
		err := f.artifacts[name].Render(&buf)
		if err != nil {
			return nil, fmt.Errorf("rendering file %q: %w", name, err)
		}
		result[name] = buf.Bytes()
	}
	return result, nil
}

// batch represents the rendered content of a complete set of artifacts, keyed
// by file name.
type batch map[string][]byte

// write stores every file under path in name order, creating path and any
// missing parents. It returns an error if a directory or a file cannot be
// written; the files stored before the failure stay in place.
func (b batch) write(path string) error {
	err := os.MkdirAll(path, 0o750)
	if err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	for _, name := range slices.Sorted(maps.Keys(b)) {
		err := b.store(filepath.Join(path, name), b[name])
		if err != nil {
			return fmt.Errorf("writing file %q: %w", name, err)
		}
	}
	return nil
}

// store writes data to path, replacing any file there and creating missing
// parent directories. It returns an error if a directory or the file cannot be
// created, written or closed.
func (b batch) store(path string, data []byte) error {
	clean := filepath.Clean(path)
	err := os.MkdirAll(filepath.Dir(clean), 0o750)
	if err != nil {
		return fmt.Errorf("creating parent directory: %w", err)
	}
	file, err := os.Create(clean)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	_, err = file.Write(data)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("writing content: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing file: %w", err)
	}
	return nil
}
