// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model/types"
)

func TestNamed_Equals(t *testing.T) {
	cases := []struct {
		name string
		a    types.Named
		b    types.Expression
		want bool
	}{
		{
			name: "returns true when both name and kind are identical",
			a:    types.NewNamed("Message", types.KindObject),
			b:    types.NewNamed("Message", types.KindObject),
			want: true,
		},
		{
			name: "returns false when the name differs but the kind is the same",
			a:    types.NewNamed("Message", types.KindObject),
			b:    types.NewNamed("Chat", types.KindObject),
			want: false,
		},
		{
			name: "returns false when the name is the same but the kind differs",
			a:    types.NewNamed("Integer", types.KindPrimitive),
			b:    types.NewNamed("Integer", types.KindObject),
			want: false,
		},
		{
			name: "returns false when compared to an Array expression",
			a:    types.NewNamed("Integer", types.KindPrimitive),
			b:    types.NewArray(types.NewNamed("Integer", types.KindPrimitive)),
			want: false,
		},
		{
			name: "returns false when compared to a Union expression",
			a:    types.NewNamed("Integer", types.KindPrimitive),
			b: types.NewUnion(
				types.NewNamed("Integer", types.KindPrimitive),
				types.NewNamed("String", types.KindPrimitive),
			),
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.a.Equals(tc.b),
				"Named.Equals must return true only when the other expression is a Named with identical name and kind",
			)
		})
	}
}

func TestNamed_String(t *testing.T) {
	cases := []struct {
		name  string
		named types.Named
		want  string
	}{
		{
			name:  "formats a primitive type as name followed by kind in parentheses",
			named: types.NewNamed("Integer", types.KindPrimitive),
			want:  "Integer(primitive)",
		},
		{
			name:  "formats an object type as name followed by kind in parentheses",
			named: types.NewNamed("Message", types.KindObject),
			want:  "Message(object)",
		},
		{
			name:  "formats a union type as name followed by kind in parentheses",
			named: types.NewNamed("ChatID", types.KindUnion),
			want:  "ChatID(union)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.named.String(),
				"Named.String must produce 'name(kind)' notation",
			)
		})
	}
}
