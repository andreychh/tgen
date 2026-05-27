// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/ir"
	"github.com/andreychh/tgen/model/types"
)

func TestField_IsInputFile(t *testing.T) {
	cases := []struct {
		name    string
		field   stubField
		want    bool
		wantErr bool
	}{
		{
			name:  "returns true when the field type name is InputFile",
			field: stubField{expr: types.NewNamed("InputFile", types.KindUnion)},
			want:  true,
		},
		{
			name:  "returns true when the field type is an Array wrapping InputFile",
			field: stubField{expr: types.NewArray(types.NewNamed("InputFile", types.KindUnion))},
			want:  true,
		},
		{
			name:  "returns false when the field type name is not InputFile",
			field: stubField{expr: types.NewNamed("Document", types.KindObject)},
			want:  false,
		},
		{
			name:  "returns false when the field type is an Array not wrapping InputFile",
			field: stubField{expr: types.NewArray(types.NewNamed("PhotoSize", types.KindObject))},
			want:  false,
		},
		{
			name:    "returns error when the field's Type method returns an error",
			field:   stubField{typeErr: errors.New("type error")},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ir.NewField(tc.field).IsInputFile()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Field.IsInputFile must propagate errors from the underlying type resolution",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Field.IsInputFile must return true only when the resolved type name is InputFile",
			)
		})
	}
}
