// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"testing"

	"github.com/andreychh/tgen/parsing"
	"github.com/stretchr/testify/assert"
)

func TestFieldKey_Value(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "accepts valid snake_case key",
			raw:  "message_id",
		},
		{
			name: "accepts key with digits",
			raw:  "chat_id_2",
		},
		{
			name:    "rejects uppercase characters",
			raw:     "MessageID",
			wantErr: true,
		},
		{
			name:    "rejects empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "rejects spaces",
			raw:     "first name",
			wantErr: true,
		},
		{
			name:    "rejects non-ASCII characters",
			raw:     "имя_поля",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parsing.NewFieldKey(tt.raw).Value()
			if tt.wantErr {
				assert.Errorf(t, err, "FieldKey must reject %q as an invalid key", tt.raw)
				return
			}
			assert.NoErrorf(t, err, "FieldKey must accept %q as a valid key", tt.raw)
		})
	}
}
