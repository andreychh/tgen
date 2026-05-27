// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/model/types"
)

func TestType_Name(t *testing.T) {
	cases := []struct {
		name    string
		expr    types.Expression
		want    string
		wantErr bool
	}{
		{
			name: "returns the name from a plain Named expression",
			expr: types.NewNamed("Integer", types.KindPrimitive),
			want: "Integer",
		},
		{
			name: "returns the inner name from a singly-nested Array expression",
			expr: types.NewArray(types.NewNamed("Message", types.KindObject)),
			want: "Message",
		},
		{
			name: "returns the inner name from a doubly-nested Array expression",
			expr: types.NewArray(types.NewArray(types.NewNamed("Chat", types.KindObject))),
			want: "Chat",
		},
		{
			name: "returns error when the expression is a Union",
			expr: types.NewUnion(
				types.NewNamed("A", types.KindPrimitive),
				types.NewNamed("B", types.KindPrimitive),
			),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ir.NewType(tc.expr).Name()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Type.Name must return an error when the expression contains a Union node",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Name must unwrap Array nesting and return the innermost Named name",
			)
		})
	}
}

func TestType_Dimensionality(t *testing.T) {
	cases := []struct {
		name    string
		expr    types.Expression
		want    int
		wantErr bool
	}{
		{
			name: "returns zero for a plain Named expression",
			expr: types.NewNamed("Integer", types.KindPrimitive),
			want: 0,
		},
		{
			name: "returns one for a singly-nested Array expression",
			expr: types.NewArray(types.NewNamed("Message", types.KindObject)),
			want: 1,
		},
		{
			name: "returns two for a doubly-nested Array expression",
			expr: types.NewArray(types.NewArray(types.NewNamed("Chat", types.KindObject))),
			want: 2,
		},
		{
			name: "returns error when the expression is a Union",
			expr: types.NewUnion(
				types.NewNamed("A", types.KindPrimitive),
				types.NewNamed("B", types.KindPrimitive),
			),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ir.NewType(tc.expr).Dimensionality()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Type.Dimensionality must return an error when the expression contains a Union node",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Dimensionality must count the number of Array wrappers around the innermost Named node",
			)
		})
	}
}

func TestType_Kind(t *testing.T) {
	cases := []struct {
		name    string
		expr    types.Expression
		want    types.Kind
		wantErr bool
	}{
		{
			name: "returns KindPrimitive from a plain Named primitive expression",
			expr: types.NewNamed("Integer", types.KindPrimitive),
			want: types.KindPrimitive,
		},
		{
			name: "returns KindObject from an Array wrapping a Named object expression",
			expr: types.NewArray(types.NewNamed("Message", types.KindObject)),
			want: types.KindObject,
		},
		{
			name: "returns KindUnion from a doubly-nested Array wrapping a Named union expression",
			expr: types.NewArray(types.NewArray(types.NewNamed("ChatID", types.KindUnion))),
			want: types.KindUnion,
		},
		{
			name: "returns error when the expression is a Union",
			expr: types.NewUnion(
				types.NewNamed("A", types.KindPrimitive),
				types.NewNamed("B", types.KindPrimitive),
			),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ir.NewType(tc.expr).Kind()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Type.Kind must return an error when the expression contains a Union node",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Kind must unwrap Array nesting and return the kind of the innermost Named node",
			)
		})
	}
}
