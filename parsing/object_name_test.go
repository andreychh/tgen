// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"testing"

	"github.com/andreychh/tgen/parsing"
	"github.com/stretchr/testify/assert"
)

func TestObjectName_Value(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "accepts PascalCase name",
			raw:  "Message",
		},
		{
			name: "accepts name with digits",
			raw:  "InputMedia3",
		},
		{
			name:    "rejects empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "rejects lowercase start",
			raw:     "message",
			wantErr: true,
		},
		{
			name:    "rejects non-ASCII characters",
			raw:     "Сообщение",
			wantErr: true,
		},
		{
			name:    "rejects underscores",
			raw:     "My_Object",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parsing.NewObjectName(tt.raw).Value()
			if tt.wantErr {
				assert.Errorf(t, err, "ObjectName must reject %q as an invalid name", tt.raw)
				return
			}
			assert.NoErrorf(t, err, "ObjectName must accept %q as a valid name", tt.raw)
		})
	}
}
