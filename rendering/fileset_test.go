// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreychh/tgen/rendering"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errRender = errors.New("simulated render failure")

type fakeView struct {
	content string
	err     error
}

func (m fakeView) Render(w io.Writer) error {
	if m.err != nil {
		return m.err
	}
	_, err := io.WriteString(w, m.content)
	return err
}

func TestFileset_Emit(t *testing.T) {
	type file struct {
		path    string
		content string
	}
	tests := []struct {
		name      string
		emitDir   string
		artifacts rendering.Artifacts
		wantFiles []file
		wantErr   bool
	}{
		{
			name:    "successfully emits multiple files and creates base directory",
			emitDir: ".",
			artifacts: rendering.Artifacts{
				"file1.go":  fakeView{content: "package api"},
				"file2.txt": fakeView{content: "test documentation"},
			},
			wantFiles: []file{
				{path: "file1.go", content: "package api"},
				{path: "file2.txt", content: "test documentation"},
			},
		},
		{
			name:    "successfully emits files into a deeply nested directory",
			emitDir: filepath.Join("foo", "bar", "api"),
			artifacts: rendering.Artifacts{
				"nested.go": fakeView{content: "package nested"},
			},
			wantFiles: []file{
				{path: filepath.Join("foo", "bar", "api", "nested.go"), content: "package nested"},
			},
		},
		{
			name:    "returns error when a view fails to render",
			emitDir: ".",
			artifacts: rendering.Artifacts{
				"bad.go": fakeView{err: errRender},
			},
			wantErr: true,
		},
		{
			name:    "returns error when file creation fails",
			emitDir: ".",
			artifacts: rendering.Artifacts{
				filepath.Join("missing_dir", "impossible.go"): fakeView{content: "this will fail"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			targetPath := filepath.Join(dir, tt.emitDir)
			err := rendering.NewFileset(tt.artifacts).Emit(targetPath)
			if tt.wantErr {
				assert.Error(t, err, "execution did not return expected error for invalid input")
				return
			}
			require.NoError(t, err, "unexpected error returned during emit execution")
			for _, wantFile := range tt.wantFiles {
				contentBytes, err := os.ReadFile(filepath.Join(dir, wantFile.path))
				require.NoError(t, err, "generated file does not read correctly")
				assert.Equal(
					t,
					wantFile.content,
					string(contentBytes),
					"file content does not match expectation",
				)
			}
		})
	}
}
