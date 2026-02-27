// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package source_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreychh/tgen/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSource_Open(t *testing.T) {
	type file struct {
		path    string
		content string
	}
	tests := []struct {
		name     string
		path     string
		file     *file
		wantData string
		wantErr  bool
	}{
		{
			name: "successfully opens existing file",
			path: "exists.txt",
			file: &file{
				path:    "exists.txt",
				content: "file content",
			},
			wantData: "file content",
		},
		{
			name:    "returns error for non-existent file",
			path:    "non-existent.txt",
			wantErr: true,
		},
		{
			name:    "returns error for directory instead of file",
			path:    "", // points to the temp directory itself
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, tt.path)
			if tt.file != nil {
				err := os.WriteFile(
					filepath.Join(dir, tt.file.path),
					[]byte(tt.file.content),
					0o644,
				)
				require.NoError(t, err, "test file was not created correctly")
			}
			rc, err := source.NewFileSource(path).Open(context.Background())
			if tt.wantErr {
				assert.Error(t, err, "did not return an error for path: %q", path)
				return
			}
			require.NoError(t, err, "did not open the file")
			data, err := io.ReadAll(rc)
			require.NoError(t, err, "did not read the file content")
			assert.Equal(t, tt.wantData, string(data), "content does not match")
			assert.NoError(t, rc.Close(), "did not close the file correctly")
		})
	}
}
