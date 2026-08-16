// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/targets/pythonv2"
)

func TestAttribute_Value(t *testing.T) {
	cases := []struct {
		name string
		key  model.Key
		want string
	}{
		{
			name: "writes a key Python reserves nothing of as it is spelled",
			key:  "chat_id",
			want: "chat_id",
		},
		{
			name: "follows a key Python reserves with an underscore",
			key:  "from",
			want: "from_",
		},
		{
			name: "leaves a soft keyword alone, which the grammar admits as a name",
			key:  "type",
			want: "type",
		},
		{
			name: "leaves the soft keyword match alone",
			key:  "match",
			want: "match",
		},
		{
			name: "follows the reserved word class with an underscore",
			key:  "class",
			want: "class_",
		},
		{
			name: "leaves a key merely opening with a reserved word alone",
			key:  "from_chat_id",
			want: "from_chat_id",
		},
		{
			name: "leaves a key differing from a reserved word by case alone",
			key:  "From",
			want: "From",
		},
		{
			name: "leaves a key already ending in an underscore alone",
			key:  "from_",
			want: "from_",
		},
		{
			name: "writes an empty key as nothing",
			key:  "",
			want: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := pythonv2.NewAttribute(c.key).Value()
			assert.Equal(
				t,
				c.want,
				got,
				"Attribute must spell a key away only where Python reserves it",
			)
		})
	}
}
