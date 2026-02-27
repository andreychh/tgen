// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package source_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreychh/tgen/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocationSource_Open_HTTP(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "successfully resolves data from URL",
			content: "http response body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(tt.content))
				},
			))
			defer server.Close()
			rc, err := source.NewLocationSource(server.URL).Open(context.Background())
			require.NoError(t, err, "did not resolve the URL location")
			data, err := io.ReadAll(rc)
			require.NoError(t, err, "did not read the content from HTTP source")
			assert.Equal(t, tt.content, string(data), "HTTP content does not match")
			assert.NoError(t, rc.Close(), "did not close HTTP stream")
		})
	}
}

func TestLocationSource_Open_File(t *testing.T) {
	type file struct {
		path    string
		content string
	}
	tests := []struct {
		name    string
		path    string
		file    *file
		content string
		wantErr bool
	}{
		{
			name:    "successfully resolves data from local file",
			path:    "test.txt",
			content: "local file content",
			file: &file{
				path:    "test.txt",
				content: "local file content",
			},
		},
		{
			name:    "returns error for non-existent file path",
			path:    "missing.txt",
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
			rc, err := source.NewLocationSource(path).Open(context.Background())
			if tt.wantErr {
				assert.Error(t, err, "did not return an error for invalid file path %q", path)
				return
			}
			require.NoError(t, err, "did not resolve the file location")
			data, err := io.ReadAll(rc)
			require.NoError(t, err, "did not read the content from file source")
			assert.Equal(t, tt.content, string(data), "file content does not match")
			assert.NoError(t, rc.Close(), "did not close file stream")
		})
	}
}
