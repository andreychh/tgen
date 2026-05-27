// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model"
	tgolang "github.com/andreychh/tgen/targets/golang"
)

func TestTag_Value(t *testing.T) {
	cases := []struct {
		name     string
		key      model.Key
		optional bool
		want     string
	}{
		{
			name:     "returns a plain JSON tag for a required field",
			key:      "message_id",
			optional: false,
			want:     "`json:\"message_id\"`",
		},
		{
			name:     "returns an omitempty JSON tag for an optional field",
			key:      "text",
			optional: true,
			want:     "`json:\"text,omitempty\"`",
		},
		{
			name:     "returns a plain JSON tag for a required field with underscores",
			key:      "chat_id",
			optional: false,
			want:     "`json:\"chat_id\"`",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tgolang.NewTag(tc.key, model.Optionality(tc.optional)).Value(),
				"Tag.Value must produce a plain JSON tag for required fields and add omitempty for optional ones",
			)
		})
	}
}
