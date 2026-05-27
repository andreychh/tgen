// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/model/types"
	tpython "github.com/andreychh/tgen/targets/python"
)

func pyAnnotation(expr types.Expression, opt bool) tpython.Annotation {
	return tpython.NewAnnotation(ir.NewType(expr), model.Optionality(opt))
}

func TestAnnotation_Value(t *testing.T) {
	cases := []struct {
		name string
		ann  tpython.Annotation
		want string
	}{
		{
			name: "returns the type annotation unchanged for a required primitive",
			ann:  pyAnnotation(types.NewNamed("Integer", types.KindPrimitive), false),
			want: "int",
		},
		{
			name: "appends None union for an optional primitive",
			ann:  pyAnnotation(types.NewNamed("Integer", types.KindPrimitive), true),
			want: "int | None = None",
		},
		{
			name: "returns the class name annotation for a required object type",
			ann:  pyAnnotation(types.NewNamed("Message", types.KindObject), false),
			want: "Message",
		},
		{
			name: "appends None union for an optional object type",
			ann:  pyAnnotation(types.NewNamed("Message", types.KindObject), true),
			want: "Message | None = None",
		},
		{
			name: "returns a list annotation for a required array type",
			ann:  pyAnnotation(types.NewArray(types.NewNamed("Message", types.KindObject)), false),
			want: "list[Message]",
		},
		{
			name: "appends None union for an optional array type",
			ann:  pyAnnotation(types.NewArray(types.NewNamed("Message", types.KindObject)), true),
			want: "list[Message] | None = None",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.ann.Value()
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Annotation.Value must return the type annotation for required fields and append ' | None = None' for optional ones",
			)
		})
	}
}
