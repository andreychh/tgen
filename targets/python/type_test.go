// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/model/types"
	tpython "github.com/andreychh/tgen/targets/python"
)

func pyType(expr types.Expression) tpython.Type {
	return tpython.NewType(ir.NewType(expr))
}

func TestType_Value(t *testing.T) {
	cases := []struct {
		name string
		typ  tpython.Type
		want string
	}{
		{
			name: "returns int for an Integer type",
			typ:  pyType(types.NewNamed("Integer", types.KindPrimitive)),
			want: "int",
		},
		{
			name: "returns float for a Float type",
			typ:  pyType(types.NewNamed("Float", types.KindPrimitive)),
			want: "float",
		},
		{
			name: "returns str for a String type",
			typ:  pyType(types.NewNamed("String", types.KindPrimitive)),
			want: "str",
		},
		{
			name: "returns bool for a Boolean type",
			typ:  pyType(types.NewNamed("Boolean", types.KindPrimitive)),
			want: "bool",
		},
		{
			name: "returns bool for a True type",
			typ:  pyType(types.NewNamed("True", types.KindPrimitive)),
			want: "bool",
		},
		{
			name: "returns the class name for an object type",
			typ:  pyType(types.NewNamed("Message", types.KindObject)),
			want: "Message",
		},
		{
			name: "returns a list annotation for a one-dimensional array of primitives",
			typ:  pyType(types.NewArray(types.NewNamed("Integer", types.KindPrimitive))),
			want: "list[int]",
		},
		{
			name: "returns a list annotation for a one-dimensional array of objects",
			typ:  pyType(types.NewArray(types.NewNamed("Message", types.KindObject))),
			want: "list[Message]",
		},
		{
			name: "returns a nested list annotation for a two-dimensional array",
			typ: pyType(
				types.NewArray(types.NewArray(types.NewNamed("Message", types.KindObject))),
			),
			want: "list[list[Message]]",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.typ.Value()
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Value must map primitives to Python built-ins and wrap arrays in list[...] for each dimension",
			)
		})
	}
}

func TestType_Part(t *testing.T) {
	cases := []struct {
		name string
		typ  tpython.Type
		want string
	}{
		{
			name: "returns the integer part template for a scalar Integer type",
			typ:  pyType(types.NewNamed("Integer", types.KindPrimitive)),
			want: "_IntPart(%s)",
		},
		{
			name: "returns the string part template for a scalar String type",
			typ:  pyType(types.NewNamed("String", types.KindPrimitive)),
			want: "_StrPart(%s)",
		},
		{
			name: "returns the bool part template for a scalar Boolean type",
			typ:  pyType(types.NewNamed("Boolean", types.KindPrimitive)),
			want: "_BoolPart(%s)",
		},
		{
			name: "returns the bare format template for a scalar object type",
			typ:  pyType(types.NewNamed("Message", types.KindObject)),
			want: "%s",
		},
		{
			name: "returns the bare format template for a scalar union type",
			typ:  pyType(types.NewNamed("ChatID", types.KindUnion)),
			want: "%s",
		},
		{
			name: "returns the object list part template for a one-dimensional array of objects",
			typ:  pyType(types.NewArray(types.NewNamed("Message", types.KindObject))),
			want: "_ObjectListPart(%s)",
		},
		{
			name: "returns the primitive list part template for a one-dimensional array of primitives",
			typ:  pyType(types.NewArray(types.NewNamed("Integer", types.KindPrimitive))),
			want: "_PrimitiveListPart(%s)",
		},
		{
			name: "returns the primitive list part template for a two-dimensional array of objects",
			typ: pyType(
				types.NewArray(types.NewArray(types.NewNamed("Message", types.KindObject))),
			),
			want: "_PrimitiveListPart(%s)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.typ.Part()
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Part must return the correct multipart form builder template for the type's name and dimensionality",
			)
		})
	}
}
