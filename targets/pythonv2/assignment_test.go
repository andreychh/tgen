// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/targets/pythonv2"
)

func TestAssignment_Value(t *testing.T) {
	cases := []struct {
		name string
		key  model.Key
		opt  model.Optionality
		want string
	}{
		{
			name: "leaves the annotation of a required field standing alone",
			key:  "chat_id",
			opt:  false,
			want: "",
		},
		{
			name: "defaults an optional field to None",
			key:  "caption",
			opt:  true,
			want: "None",
		},
		{
			name: "keeps the key of a required field Python reserves the name of",
			key:  "from",
			opt:  false,
			want: `Field(alias="from")`,
		},
		{
			name: "keeps the key beside the default where such a field is optional",
			key:  "from",
			opt:  true,
			want: `Field(default=None, alias="from")`,
		},
		{
			name: "writes no alias for a key merely opening with a reserved word",
			key:  "from_chat_id",
			opt:  true,
			want: "None",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := pythonv2.NewAssignment(c.key, c.opt).Value()
			assert.Equal(
				t,
				c.want,
				got,
				"Assignment must name the key only where the attribute stops spelling it",
			)
		})
	}
}
