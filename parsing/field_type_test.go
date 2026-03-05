// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"testing"

	"github.com/andreychh/tgen/parsing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldType_Value(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{
			name: "accepts simple named type",
			raw:  "Integer",
			want: "Integer",
		},
		{
			name: "accepts array type",
			raw:  "Array of String",
			want: "Array of String",
		},
		{
			name: "accepts union type",
			raw:  "Integer or String",
			want: "Integer or String",
		},
		{
			name:    "rejects empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "rejects non-ASCII characters",
			raw:     "Целое",
			wantErr: true,
		},
		{
			name:    "rejects special characters",
			raw:     "Array<Integer>",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsing.NewFieldType(tt.raw).Value()
			if tt.wantErr {
				assert.Errorf(t, err, "FieldType must reject %q as an invalid type", tt.raw)
				return
			}
			require.NoErrorf(t, err, "FieldType must accept %q as a valid type", tt.raw)
			assert.Equalf(
				t,
				tt.want,
				got,
				"FieldType must return the original string for %q",
				tt.raw,
			)
		})
	}
}
