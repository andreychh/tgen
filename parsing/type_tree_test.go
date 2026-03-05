// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"testing"

	"github.com/andreychh/tgen/parsing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type typeAssertion func(t *testing.T, expr parsing.TypeExpression)

func TestTypeTree_Root(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		verify  typeAssertion
		wantErr bool
	}{
		{
			name:   "parses named type",
			raw:    "Integer",
			verify: assertTypeTree_Named("Integer"),
		},
		{
			name:   "parses array of named type",
			raw:    "Array of String",
			verify: assertTypeTree_Array("String"),
		},
		{
			name:   "parses nested array",
			raw:    "Array of Array of Integer",
			verify: assertTypeTree_NestedArray("Integer"),
		},
		{
			name:   "parses union with or",
			raw:    "Integer or String",
			verify: assertTypeTree_Union(2),
		},
		{
			name:   "parses union with comma-and syntax",
			raw:    "InputMediaAudio, InputMediaDocument and InputMediaVideo",
			verify: assertTypeTree_Union(3),
		},
		{
			name:   "parses array of union",
			raw:    "Array of InputMediaAudio, InputMediaDocument and InputMediaVideo",
			verify: assertTypeTree_ArrayOfUnion(3),
		},
		{
			name:    "rejects empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "rejects invalid characters",
			raw:     "Array<Integer>",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parsing.NewTypeTree(parsing.NewFieldType(tt.raw)).Root()
			if tt.wantErr {
				assert.Errorf(t, err, "TypeTree must reject %q as invalid input", tt.raw)
				return
			}
			require.NoErrorf(t, err, "TypeTree must parse %q without error", tt.raw)
			tt.verify(t, expr)
		})
	}
}

func assertTypeTree_Named(want string) typeAssertion {
	return func(t *testing.T, expr parsing.TypeExpression) {
		t.Helper()
		name, ok := expr.Named()
		require.Truef(t, ok, "TypeTree must produce a named node")
		assert.Equalf(t, want, name, "TypeTree must preserve type name %q", want)
	}
}

func assertTypeTree_Array(wantInner string) typeAssertion {
	return func(t *testing.T, expr parsing.TypeExpression) {
		t.Helper()
		inner, ok := expr.Array()
		require.Truef(t, ok, "TypeTree must produce an array node")
		assertTypeTree_Named(wantInner)(t, inner)
	}
}

func assertTypeTree_NestedArray(wantInnermost string) typeAssertion {
	return func(t *testing.T, expr parsing.TypeExpression) {
		t.Helper()
		inner, ok := expr.Array()
		require.Truef(t, ok, "TypeTree must produce outer array node")
		assertTypeTree_Array(wantInnermost)(t, inner)
	}
}

func assertTypeTree_Union(wantVariants int) typeAssertion {
	return func(t *testing.T, expr parsing.TypeExpression) {
		t.Helper()
		variants, ok := expr.Union()
		require.Truef(t, ok, "TypeTree must produce a union node")
		assert.Lenf(t, variants, wantVariants, "TypeTree must produce %d variants", wantVariants)
	}
}

func assertTypeTree_ArrayOfUnion(wantVariants int) typeAssertion {
	return func(t *testing.T, expr parsing.TypeExpression) {
		t.Helper()
		inner, ok := expr.Array()
		require.Truef(t, ok, "TypeTree must produce array node wrapping union")
		assertTypeTree_Union(wantVariants)(t, inner)
	}
}
