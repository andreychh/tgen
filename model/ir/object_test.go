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

func TestObject_HasInputFile(t *testing.T) {
	inputFile := stubField{expr: types.NewNamed("InputFile", types.KindUnion)}
	document := stubField{expr: types.NewNamed("Document", types.KindObject)}
	cases := []struct {
		name    string
		fields  []stubField
		want    bool
		wantErr bool
	}{
		{
			name:   "returns false when the object has no fields",
			fields: nil,
			want:   false,
		},
		{
			name:   "returns false when no field has the InputFile type",
			fields: []stubField{document},
			want:   false,
		},
		{
			name:   "returns true when the only field has the InputFile type",
			fields: []stubField{inputFile},
			want:   true,
		},
		{
			name:   "returns true when InputFile is among multiple fields",
			fields: []stubField{document, inputFile},
			want:   true,
		},
		{
			name:    "returns error when a field's Type method returns an error",
			fields:  []stubField{{typeErr: errors.New("type error")}},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			specFields := make([]stubField, len(tc.fields))
			copy(specFields, tc.fields)
			obj := stubObject{fields: toSpecFields(specFields)}
			got, err := ir.NewObject(obj).HasInputFile()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Object.HasInputFile must propagate errors from field type resolution",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Object.HasInputFile must return true only when at least one field resolves to the InputFile type",
			)
		})
	}
}
