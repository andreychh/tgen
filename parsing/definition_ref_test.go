// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"testing"

	"github.com/andreychh/tgen/parsing"
	"github.com/stretchr/testify/assert"
)

func TestDefinitionRef_Value(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "accepts lowercase letters",
			raw:  "message",
		},
		{
			name: "accepts letters with digits",
			raw:  "inputmedia3",
		},
		{
			name:    "rejects empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "rejects uppercase characters",
			raw:     "Message",
			wantErr: true,
		},
		{
			name:    "rejects hash prefix",
			raw:     "#message",
			wantErr: true,
		},
		{
			name:    "rejects non-ASCII characters",
			raw:     "сообщение",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parsing.NewDefinitionRef(tt.raw).Value()
			if tt.wantErr {
				assert.Errorf(t, err, "DefinitionRef must reject %q as an invalid ref", tt.raw)
				return
			}
			assert.NoErrorf(t, err, "DefinitionRef must accept %q as a valid ref", tt.raw)
		})
	}
}
