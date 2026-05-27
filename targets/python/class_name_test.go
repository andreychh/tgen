// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model"
	tpython "github.com/andreychh/tgen/targets/python"
)

func TestClassName_Value(t *testing.T) {
	cases := []struct {
		name  string
		input model.Name
		want  string
	}{
		{
			name:  "converts a snake_case name to PascalCase",
			input: "message_text",
			want:  "MessageText",
		},
		{
			name:  "replaces Id suffix with ID",
			input: "chat_id",
			want:  "ChatID",
		},
		{
			name:  "replaces Url with URL",
			input: "callback_url",
			want:  "CallbackURL",
		},
		{
			name:  "replaces Api with API",
			input: "api_hash",
			want:  "APIHash",
		},
		{
			name:  "replaces Ip with IP",
			input: "local_ip",
			want:  "LocalIP",
		},
		{
			name:  "returns a single-word PascalCase name unchanged",
			input: "Message",
			want:  "Message",
		},
		{
			name:  "converts a single lowercase word to PascalCase",
			input: "update",
			want:  "Update",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tpython.NewClassName(tc.input).Value(),
				"ClassName.Value must produce a PascalCase Python class name with standard acronym replacements",
			)
		})
	}
}
