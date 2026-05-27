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

func TestDiscriminatedObject_HasInputFile(t *testing.T) {
	inputFile := stubField{expr: types.NewNamed("InputFile", types.KindUnion)}
	thumbnail := stubField{expr: types.NewNamed("PhotoSize", types.KindObject)}
	cases := []struct {
		name    string
		fields  []stubField
		want    bool
		wantErr bool
	}{
		{
			name:   "returns false when the variant has no free fields",
			fields: nil,
			want:   false,
		},
		{
			name:   "returns false when no free field has the InputFile type",
			fields: []stubField{thumbnail},
			want:   false,
		},
		{
			name:   "returns true when a free field has the InputFile type",
			fields: []stubField{thumbnail, inputFile},
			want:   true,
		},
		{
			name:    "returns error when a free field's Type method returns an error",
			fields:  []stubField{{typeErr: errors.New("type error")}},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obj := stubDiscriminatedObject{fields: toSpecFields(tc.fields)}
			got, err := ir.NewDiscriminatedObject(obj).HasInputFile()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"DiscriminatedObject.HasInputFile must propagate errors from free field type resolution",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"DiscriminatedObject.HasInputFile must return true only when at least one free field resolves to the InputFile type",
			)
		})
	}
}
