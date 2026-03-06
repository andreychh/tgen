// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"testing"

	"github.com/andreychh/tgen/parsing"
	"github.com/stretchr/testify/assert"
)

func TestReleaseRef_Value(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "accepts valid release ref",
			raw:  "february-9-2026",
		},
		{
			name: "accepts multi-word month",
			raw:  "december-27-2024",
		},
		{
			name:    "rejects empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "rejects hash prefix",
			raw:     "#february-9-2026",
			wantErr: true,
		},
		{
			name:    "rejects missing date",
			raw:     "february",
			wantErr: true,
		},
		{
			name:    "rejects non-ASCII characters",
			raw:     "февраль-9-2026",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parsing.NewReleaseRef(tt.raw).Value()
			if tt.wantErr {
				assert.Errorf(t, err, "ReleaseRef must reject %q as an invalid ref", tt.raw)
				return
			}
			assert.NoErrorf(t, err, "ReleaseRef must accept %q as a valid ref", tt.raw)
		})
	}
}
