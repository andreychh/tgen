// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/spec/overlays"
	"github.com/andreychh/tgen/model/types"
)

func TestInputFile_Apply(t *testing.T) {
	inputFileObj := types.NewNamed("InputFile", types.KindObject)
	inputFileUnion := types.NewNamed("InputFile", types.KindUnion)
	str := types.NewNamed("String", types.KindPrimitive)
	cases := []struct {
		name    string
		field   spec.Field
		want    types.Expression
		wantErr bool
	}{
		{
			name:  "replaces bare InputFile object type with InputFile union type",
			field: stubField{expr: inputFileObj, desc: stubDesc{}},
			want:  inputFileUnion,
		},
		{
			name:  "replaces InputFile or String union with InputFile union type",
			field: stubField{expr: types.NewUnion(inputFileObj, str), desc: stubDesc{}},
			want:  inputFileUnion,
		},
		{
			name:  "replaces String type that links to sending-files with InputFile union type",
			field: stubField{expr: str, desc: stubDesc{links: []string{"#sending-files"}}},
			want:  inputFileUnion,
		},
		{
			name:  "passes through String type that has no sending-files link",
			field: stubField{expr: str, desc: stubDesc{links: []string{"#other-section"}}},
			want:  str,
		},
		{
			name:  "passes through String type when there are no description links at all",
			field: stubField{expr: str, desc: stubDesc{}},
			want:  str,
		},
		{
			name:  "passes through a non-matching named type",
			field: stubField{expr: types.NewNamed("Document", types.KindObject), desc: stubDesc{}},
			want:  types.NewNamed("Document", types.KindObject),
		},
		{
			name:    "passes through unchanged when Type returns an error",
			field:   stubField{typeErr: errType, desc: stubDesc{}},
			wantErr: true,
		},
		{
			name:  "passes through unchanged when Description Links returns an error",
			field: stubField{expr: str, desc: stubDesc{linksErr: errLinks}},
			want:  str,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlays.InputFile{}.Apply(tc.field).Type()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"InputFile.Apply must return the original field unchanged when its Type method errors",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"InputFile.Apply must replace InputFile-bearing field types with the InputFile union and pass through all others",
			)
		})
	}
}
