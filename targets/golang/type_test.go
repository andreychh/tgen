// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/model/types"
	tgolang "github.com/andreychh/tgen/targets/golang"
)

func makeType(expr types.Expression, opt bool) tgolang.Type {
	return tgolang.NewType(ir.NewType(expr), model.Optionality(opt))
}

func TestType_Shape(t *testing.T) {
	cases := []struct {
		name string
		typ  tgolang.Type
		want tgolang.Shape
	}{
		{
			name: "returns ShapePlain for a required primitive named type",
			typ:  makeType(types.NewNamed("Integer", types.KindPrimitive), false),
			want: tgolang.ShapePlain,
		},
		{
			name: "returns ShapePlain for a required array of object type",
			typ:  makeType(types.NewArray(types.NewNamed("Message", types.KindObject)), false),
			want: tgolang.ShapePlain,
		},
		{
			name: "returns ShapeUnion for a required scalar union type",
			typ:  makeType(types.NewNamed("ChatID", types.KindUnion), false),
			want: tgolang.ShapeUnion,
		},
		{
			name: "returns ShapeUnionArray for a required one-dimensional union array",
			typ:  makeType(types.NewArray(types.NewNamed("ChatID", types.KindUnion)), false),
			want: tgolang.ShapeUnionArray,
		},
		{
			name: "returns ShapeUnion for a two-dimensional union array",
			typ: makeType(
				types.NewArray(types.NewArray(types.NewNamed("ChatID", types.KindUnion))),
				false,
			),
			want: tgolang.ShapeUnion,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.typ.Shape()
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Shape must return ShapeUnionArray for a one-dimensional union, ShapeUnion for other unions, and ShapePlain for everything else",
			)
		})
	}
}

func TestType_Value(t *testing.T) {
	cases := []struct {
		name string
		typ  tgolang.Type
		want string
	}{
		{
			name: "returns int64 for a required Integer type",
			typ:  makeType(types.NewNamed("Integer", types.KindPrimitive), false),
			want: "int64",
		},
		{
			name: "returns float64 for a required Float type",
			typ:  makeType(types.NewNamed("Float", types.KindPrimitive), false),
			want: "float64",
		},
		{
			name: "returns string for a required String type",
			typ:  makeType(types.NewNamed("String", types.KindPrimitive), false),
			want: "string",
		},
		{
			name: "returns bool for a required Boolean type",
			typ:  makeType(types.NewNamed("Boolean", types.KindPrimitive), false),
			want: "bool",
		},
		{
			name: "returns bool for a required True type",
			typ:  makeType(types.NewNamed("True", types.KindPrimitive), false),
			want: "bool",
		},
		{
			name: "returns the struct name for a required object type",
			typ:  makeType(types.NewNamed("Message", types.KindObject), false),
			want: "Message",
		},
		{
			name: "returns a pointer to the struct for an optional scalar object type",
			typ:  makeType(types.NewNamed("Message", types.KindObject), true),
			want: "*Message",
		},
		{
			name: "returns a slice for a required array of primitives",
			typ:  makeType(types.NewArray(types.NewNamed("Integer", types.KindPrimitive)), false),
			want: "[]int64",
		},
		{
			name: "returns a slice for a required array of objects",
			typ:  makeType(types.NewArray(types.NewNamed("Message", types.KindObject)), false),
			want: "[]Message",
		},
		{
			name: "returns a slice without a pointer for an optional array type",
			typ:  makeType(types.NewArray(types.NewNamed("Message", types.KindObject)), true),
			want: "[]Message",
		},
		{
			name: "returns the union name without a pointer for an optional scalar union type",
			typ:  makeType(types.NewNamed("ChatID", types.KindUnion), true),
			want: "ChatID",
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
				"Type.Value must map primitives to Go built-ins, add a pointer for optional scalar non-union types, and prefix arrays with []",
			)
		})
	}
}

func TestType_Zero(t *testing.T) {
	cases := []struct {
		name string
		typ  tgolang.Type
		want string
	}{
		{
			name: "returns nil for an optional type",
			typ:  makeType(types.NewNamed("Message", types.KindObject), true),
			want: "nil",
		},
		{
			name: "returns nil for a required array type",
			typ:  makeType(types.NewArray(types.NewNamed("Integer", types.KindPrimitive)), false),
			want: "nil",
		},
		{
			name: "returns nil for a required scalar union type",
			typ:  makeType(types.NewNamed("ChatID", types.KindUnion), false),
			want: "nil",
		},
		{
			name: "returns 0 for a required Integer type",
			typ:  makeType(types.NewNamed("Integer", types.KindPrimitive), false),
			want: "0",
		},
		{
			name: "returns an empty string literal for a required String type",
			typ:  makeType(types.NewNamed("String", types.KindPrimitive), false),
			want: `""`,
		},
		{
			name: "returns false for a required Boolean type",
			typ:  makeType(types.NewNamed("Boolean", types.KindPrimitive), false),
			want: "false",
		},
		{
			name: "returns a struct literal zero for a required object type",
			typ:  makeType(types.NewNamed("Message", types.KindObject), false),
			want: "Message{}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.typ.Zero()
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Zero must return nil for optional, array, and union types; primitive literals for primitives; and struct{} for objects",
			)
		})
	}
}

func TestType_Part(t *testing.T) {
	cases := []struct {
		name string
		typ  tgolang.Type
		want string
	}{
		{
			name: "returns newPrimitiveSlicePart template for a required array of primitives",
			typ:  makeType(types.NewArray(types.NewNamed("Integer", types.KindPrimitive)), false),
			want: "newPrimitiveSlicePart(%s)",
		},
		{
			name: "returns newObjectSlicePart template for a required array of objects",
			typ:  makeType(types.NewArray(types.NewNamed("Message", types.KindObject)), false),
			want: "newObjectSlicePart(%s)",
		},
		{
			name: "returns bare format template for a required scalar object type",
			typ:  makeType(types.NewNamed("Message", types.KindObject), false),
			want: "%s",
		},
		{
			name: "returns bare format template for a required scalar union type",
			typ:  makeType(types.NewNamed("ChatID", types.KindUnion), false),
			want: "%s",
		},
		{
			name: "returns primitive part template with direct argument for a required scalar primitive",
			typ:  makeType(types.NewNamed("Integer", types.KindPrimitive), false),
			want: "newInt64Part(%s)",
		},
		{
			name: "returns primitive part template with pointer dereference for an optional scalar primitive",
			typ:  makeType(types.NewNamed("Integer", types.KindPrimitive), true),
			want: "newInt64Part(*%s)",
		},
		{
			name: "returns string part template with pointer dereference for an optional scalar string",
			typ:  makeType(types.NewNamed("String", types.KindPrimitive), true),
			want: "newStringPart(*%s)",
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
				"Type.Part must return the correct multipart form builder template string for the field type and optionality",
			)
		})
	}
}
