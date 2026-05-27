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

func TestChatID_Apply(t *testing.T) {
	integer := types.NewNamed("Integer", types.KindPrimitive)
	str := types.NewNamed("String", types.KindPrimitive)
	cases := []struct {
		name    string
		field   spec.Field
		want    types.Expression
		wantErr bool
	}{
		{
			name:  "replaces Integer or String union with ChatID named type",
			field: stubField{expr: types.NewUnion(integer, str)},
			want:  types.NewNamed("ChatID", types.KindUnion),
		},
		{
			name:  "passes through a field whose type is neither Integer nor String",
			field: stubField{expr: types.NewNamed("Boolean", types.KindPrimitive)},
			want:  types.NewNamed("Boolean", types.KindPrimitive),
		},
		{
			name:  "passes through Integer alone without replacement",
			field: stubField{expr: integer},
			want:  integer,
		},
		{
			name:    "passes through unchanged when Type returns an error",
			field:   stubField{typeErr: errType},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlays.ChatID{}.Apply(tc.field).Type()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ChatID.Apply must return the original field unchanged when its Type method errors",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ChatID.Apply must replace Integer|String with ChatID and pass through all other types",
			)
		})
	}
}
