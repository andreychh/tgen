// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model/types"
)

func TestUnion_Equals(t *testing.T) {
	integer := types.NewNamed("Integer", types.KindPrimitive)
	str := types.NewNamed("String", types.KindPrimitive)
	msg := types.NewNamed("Message", types.KindObject)
	cases := []struct {
		name string
		a    types.Union
		b    types.Expression
		want bool
	}{
		{
			name: "returns true when all variants match in the same order",
			a:    types.NewUnion(integer, str),
			b:    types.NewUnion(integer, str),
			want: true,
		},
		{
			name: "returns false when variants are in a different order",
			a:    types.NewUnion(integer, str),
			b:    types.NewUnion(str, integer),
			want: false,
		},
		{
			name: "returns false when one variant name differs",
			a:    types.NewUnion(integer, str),
			b:    types.NewUnion(integer, msg),
			want: false,
		},
		{
			name: "returns false when the variant count differs",
			a:    types.NewUnion(integer, str),
			b:    types.NewUnion(integer, str, msg),
			want: false,
		},
		{
			name: "returns false when compared to a Named expression",
			a:    types.NewUnion(integer, str),
			b:    integer,
			want: false,
		},
		{
			name: "returns false when compared to an Array expression",
			a:    types.NewUnion(integer, str),
			b:    types.NewArray(integer),
			want: false,
		},
		{
			name: "returns true for a single-element union compared to itself",
			a:    types.NewUnion(integer),
			b:    types.NewUnion(integer),
			want: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.a.Equals(tc.b),
				"Union.Equals must return true only when the other expression is a Union with identical variants in the same order",
			)
		})
	}
}

func TestUnion_String(t *testing.T) {
	cases := []struct {
		name  string
		union types.Union
		want  string
	}{
		{
			name: "joins variant strings with pipe separators",
			union: types.NewUnion(
				types.NewNamed("Integer", types.KindPrimitive),
				types.NewNamed("String", types.KindPrimitive),
			),
			want: "Integer(primitive) | String(primitive)",
		},
		{
			name:  "returns the single variant string when the union has one element",
			union: types.NewUnion(types.NewNamed("Message", types.KindObject)),
			want:  "Message(object)",
		},
		{
			name: "joins three variants with pipe separators in declaration order",
			union: types.NewUnion(
				types.NewNamed("A", types.KindPrimitive),
				types.NewNamed("B", types.KindPrimitive),
				types.NewNamed("C", types.KindPrimitive),
			),
			want: "A(primitive) | B(primitive) | C(primitive)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.union.String(),
				"Union.String must join all variant strings with ' | ' separators in declaration order",
			)
		})
	}
}
