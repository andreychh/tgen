// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"testing"

	"github.com/andreychh/tgen/parsing"
	"github.com/stretchr/testify/assert"
)

func TestMethodName_Value(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "accepts camelCase name",
			raw:  "sendMessage",
		},
		{
			name: "accepts name with digits",
			raw:  "getUpdates2",
		},
		{
			name:    "rejects empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "rejects PascalCase start",
			raw:     "SendMessage",
			wantErr: true,
		},
		{
			name:    "rejects non-ASCII characters",
			raw:     "отправитьСообщение",
			wantErr: true,
		},
		{
			name:    "rejects underscores",
			raw:     "send_message",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parsing.NewMethodName(tt.raw).Value()
			if tt.wantErr {
				assert.Errorf(t, err, "MethodName must reject %q as an invalid name", tt.raw)
				return
			}
			assert.NoErrorf(t, err, "MethodName must accept %q as a valid name", tt.raw)
		})
	}
}
