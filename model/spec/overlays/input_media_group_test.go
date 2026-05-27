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

func inputMediaGroupArray() types.Expression {
	return types.NewArray(types.NewUnion(
		types.NewNamed("InputMediaAudio", types.KindObject),
		types.NewNamed("InputMediaDocument", types.KindObject),
		types.NewNamed("InputMediaLivePhoto", types.KindObject),
		types.NewNamed("InputMediaPhoto", types.KindObject),
		types.NewNamed("InputMediaVideo", types.KindObject),
	))
}

func TestInputMediaGroup_Apply(t *testing.T) {
	cases := []struct {
		name    string
		field   spec.Field
		want    types.Expression
		wantErr bool
	}{
		{
			name:  "replaces the five-variant input media array union with InputMediaGroup array",
			field: stubField{expr: inputMediaGroupArray()},
			want:  types.NewArray(types.NewNamed("InputMediaGroup", types.KindUnion)),
		},
		{
			name: "passes through a four-variant input media array that is missing one variant",
			field: stubField{expr: types.NewArray(types.NewUnion(
				types.NewNamed("InputMediaAudio", types.KindObject),
				types.NewNamed("InputMediaDocument", types.KindObject),
				types.NewNamed("InputMediaPhoto", types.KindObject),
				types.NewNamed("InputMediaVideo", types.KindObject),
			))},
			want: types.NewArray(types.NewUnion(
				types.NewNamed("InputMediaAudio", types.KindObject),
				types.NewNamed("InputMediaDocument", types.KindObject),
				types.NewNamed("InputMediaPhoto", types.KindObject),
				types.NewNamed("InputMediaVideo", types.KindObject),
			)),
		},
		{
			name:  "passes through a plain named array type without replacement",
			field: stubField{expr: types.NewArray(types.NewNamed("Message", types.KindObject))},
			want:  types.NewArray(types.NewNamed("Message", types.KindObject)),
		},
		{
			name:    "passes through unchanged when Type returns an error",
			field:   stubField{typeErr: errType},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlays.InputMediaGroup{}.Apply(tc.field).Type()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"InputMediaGroup.Apply must return the original field unchanged when its Type method errors",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"InputMediaGroup.Apply must replace the five-variant input media array with InputMediaGroup array and pass through all other types",
			)
		})
	}
}
