// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/types"
)

type stubCatalog map[string]types.Kind

func (c stubCatalog) Lookup(name string) (types.Kind, bool) {
	k, ok := c[name]
	return k, ok
}

type stubSource struct {
	value string
	err   error
}

func (s stubSource) Value() (string, error) { return s.value, s.err }

var testCatalog = stubCatalog{
	"Integer": types.KindPrimitive,
	"String":  types.KindPrimitive,
	"Boolean": types.KindPrimitive,
	"Message": types.KindObject,
	"Chat":    types.KindObject,
	"ChatID":  types.KindUnion,
}

func src(value string) stubSource { return stubSource{value: value} }
func srcErr(err error) stubSource { return stubSource{err: err} }

func TestType_Value(t *testing.T) {
	integer := types.NewNamed("Integer", types.KindPrimitive)
	str := types.NewNamed("String", types.KindPrimitive)
	message := types.NewNamed("Message", types.KindObject)
	chat := types.NewNamed("Chat", types.KindObject)
	cases := []struct {
		name    string
		source  stubSource
		want    types.Expression
		wantErr bool
	}{
		{
			name:   "parses a single known type name into a Named expression",
			source: src("Integer"),
			want:   integer,
		},
		{
			name:   "parses Array of X into an Array wrapping the element type",
			source: src("Array of Integer"),
			want:   types.NewArray(integer),
		},
		{
			name:   "parses A or B into a two-variant Union",
			source: src("Integer or String"),
			want:   types.NewUnion(integer, str),
		},
		{
			name:   "parses comma-separated and-joined list into a three-variant Union",
			source: src("Integer, String and Boolean"),
			want:   types.NewUnion(integer, str, types.NewNamed("Boolean", types.KindPrimitive)),
		},
		{
			name:   "parses nested Array of Array of X into a doubly-nested Array",
			source: src("Array of Array of Integer"),
			want:   types.NewArray(types.NewArray(integer)),
		},
		{
			name:   "parses Array of A or B into an Array wrapping a two-variant Union",
			source: src("Array of Integer or String"),
			want:   types.NewArray(types.NewUnion(integer, str)),
		},
		{
			name:   "parses a two-variant union of object types",
			source: src("Message or Chat"),
			want:   types.NewUnion(message, chat),
		},
		{
			name:    "returns error when the type source returns an error",
			source:  srcErr(errors.New("source error")),
			wantErr: true,
		},
		{
			name:    "returns error when the expression string is empty",
			source:  src(""),
			wantErr: true,
		},
		{
			name:    "returns error for Array of with no element name",
			source:  src("Array of "),
			wantErr: true,
		},
		{
			name:    "returns error when the type name is not in the catalog",
			source:  src("UnknownType"),
			wantErr: true,
		},
		{
			name:    "returns error when a union variant is not in the catalog",
			source:  src("Integer or UnknownType"),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := types.NewType(testCatalog, tc.source).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Type.Value must return an error when the source errors or the expression cannot be parsed",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Type.Value must parse the source string into the correct Expression tree",
			)
		})
	}
}
