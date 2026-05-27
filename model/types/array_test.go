// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model/types"
)

func TestArray_Equals(t *testing.T) {
	integer := types.NewNamed("Integer", types.KindPrimitive)
	str := types.NewNamed("String", types.KindPrimitive)
	cases := []struct {
		name string
		a    types.Array
		b    types.Expression
		want bool
	}{
		{
			name: "returns true when both arrays wrap the same element type",
			a:    types.NewArray(integer),
			b:    types.NewArray(integer),
			want: true,
		},
		{
			name: "returns false when the element types differ by name",
			a:    types.NewArray(integer),
			b:    types.NewArray(str),
			want: false,
		},
		{
			name: "returns false when the element types differ by kind",
			a:    types.NewArray(types.NewNamed("Integer", types.KindPrimitive)),
			b:    types.NewArray(types.NewNamed("Integer", types.KindObject)),
			want: false,
		},
		{
			name: "returns false when compared to a Named expression",
			a:    types.NewArray(integer),
			b:    integer,
			want: false,
		},
		{
			name: "returns false when compared to a Union expression",
			a:    types.NewArray(integer),
			b:    types.NewUnion(integer, str),
			want: false,
		},
		{
			name: "returns true for nested arrays with identical element chains",
			a:    types.NewArray(types.NewArray(integer)),
			b:    types.NewArray(types.NewArray(integer)),
			want: true,
		},
		{
			name: "returns false for nested arrays when inner elements differ",
			a:    types.NewArray(types.NewArray(integer)),
			b:    types.NewArray(types.NewArray(str)),
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.a.Equals(tc.b),
				"Array.Equals must return true only when the other expression is an Array with an equal element type",
			)
		})
	}
}

func TestArray_String(t *testing.T) {
	cases := []struct {
		name  string
		array types.Array
		want  string
	}{
		{
			name:  "formats a named element type in Array<element> notation",
			array: types.NewArray(types.NewNamed("Integer", types.KindPrimitive)),
			want:  "Array<Integer(primitive)>",
		},
		{
			name:  "formats doubly-nested arrays in Array<Array<element>> notation",
			array: types.NewArray(types.NewArray(types.NewNamed("String", types.KindPrimitive))),
			want:  "Array<Array<String(primitive)>>",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.array.String(),
				"Array.String must produce 'Array<element>' notation recursively",
			)
		})
	}
}
