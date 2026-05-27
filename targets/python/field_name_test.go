// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model"
	tpython "github.com/andreychh/tgen/targets/python"
)

func TestFieldName_Value(t *testing.T) {
	cases := []struct {
		name string
		key  model.Key
		want string
	}{
		{
			name: "returns the key unchanged for a plain snake_case field",
			key:  "message_id",
			want: "message_id",
		},
		{
			name: "returns the key unchanged for a single-word field",
			key:  "text",
			want: "text",
		},
		{
			name: "appends underscore to the Python keyword from",
			key:  "from",
			want: "from_",
		},
		{
			name: "appends underscore to the Python keyword import",
			key:  "import",
			want: "import_",
		},
		{
			name: "appends underscore to the Python keyword class",
			key:  "class",
			want: "class_",
		},
		{
			name: "appends underscore to the Python keyword return",
			key:  "return",
			want: "return_",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tpython.NewFieldName(tc.key).Value(),
				"FieldName.Value must return the key unchanged except for Python keywords, which get a trailing underscore",
			)
		})
	}
}
