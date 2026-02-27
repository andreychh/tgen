// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package source_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreychh/tgen/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPSource_Open(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		content    string
		wantErr    bool
	}{
		{
			name:       "successfully fetches data from remote url",
			statusCode: http.StatusOK,
			content:    "remote content",
		},
		{
			name:       "returns error for non-200 status code",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "returns error for internal server error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write([]byte(tt.content))
				},
			))
			defer server.Close()
			rc, err := source.NewDefaultHTTPSource(server.URL).Open(context.Background())
			if tt.wantErr {
				assert.Error(t, err, "did not return an error for status code: %d", tt.statusCode)
				return
			}
			require.NoError(t, err, "did not fetch the data correctly")
			data, err := io.ReadAll(rc)
			require.NoError(t, err, "did not read the response body")
			assert.Equal(t, tt.content, string(data), "remote content does not match")
			assert.NoError(t, rc.Close(), "did not close the response body correctly")
		})
	}
}

func TestHTTPSource_Open_ContextCancellation(t *testing.T) {
	t.Run(
		"returns error when context is cancelled",
		func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			))
			defer server.Close()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := source.NewDefaultHTTPSource(server.URL).Open(ctx)
			assert.Error(t, err, "did not return an error upon context cancellation")
		},
	)
}
