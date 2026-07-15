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

func inputRichMediaUnion() types.Expression {
	return types.NewUnion(
		types.NewNamed("InputMediaAnimation", types.KindObject),
		types.NewNamed("InputMediaAudio", types.KindObject),
		types.NewNamed("InputMediaPhoto", types.KindObject),
		types.NewNamed("InputMediaVideo", types.KindObject),
		types.NewNamed("InputMediaVoiceNote", types.KindObject),
	)
}

func TestInputRichMedia_Apply(t *testing.T) {
	cases := []struct {
		name    string
		field   spec.Field
		want    types.Expression
		wantErr bool
	}{
		{
			name:  "replaces the five-variant input rich media union with InputRichMedia",
			field: stubField{expr: inputRichMediaUnion()},
			want:  types.NewNamed("InputRichMedia", types.KindUnion),
		},
		{
			name: "passes through a four-variant union that is missing one variant",
			field: stubField{expr: types.NewUnion(
				types.NewNamed("InputMediaAnimation", types.KindObject),
				types.NewNamed("InputMediaAudio", types.KindObject),
				types.NewNamed("InputMediaPhoto", types.KindObject),
				types.NewNamed("InputMediaVideo", types.KindObject),
			)},
			want: types.NewUnion(
				types.NewNamed("InputMediaAnimation", types.KindObject),
				types.NewNamed("InputMediaAudio", types.KindObject),
				types.NewNamed("InputMediaPhoto", types.KindObject),
				types.NewNamed("InputMediaVideo", types.KindObject),
			),
		},
		{
			name:  "passes through a plain named type without replacement",
			field: stubField{expr: types.NewNamed("Message", types.KindObject)},
			want:  types.NewNamed("Message", types.KindObject),
		},
		{
			name:    "passes through unchanged when Type returns an error",
			field:   stubField{typeErr: errType},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlays.InputRichMedia{}.Apply(tc.field).Type()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"InputRichMedia.Apply must return the original field unchanged when its Type method errors",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"InputRichMedia.Apply must replace the five-variant input rich media union with InputRichMedia and pass through all other types",
			)
		})
	}
}
