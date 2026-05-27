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

func replyMarkupUnion() types.Expression {
	return types.NewUnion(
		types.NewNamed("InlineKeyboardMarkup", types.KindObject),
		types.NewNamed("ReplyKeyboardMarkup", types.KindObject),
		types.NewNamed("ReplyKeyboardRemove", types.KindObject),
		types.NewNamed("ForceReply", types.KindObject),
	)
}

func TestReplyMarkup_Apply(t *testing.T) {
	cases := []struct {
		name    string
		field   spec.Field
		want    types.Expression
		wantErr bool
	}{
		{
			name:  "replaces the four-variant reply markup union with ReplyMarkup named type",
			field: stubField{expr: replyMarkupUnion()},
			want:  types.NewNamed("ReplyMarkup", types.KindUnion),
		},
		{
			name: "passes through a three-variant union that is not the reply markup union",
			field: stubField{
				expr: types.NewUnion(
					types.NewNamed("InlineKeyboardMarkup", types.KindObject),
					types.NewNamed("ReplyKeyboardMarkup", types.KindObject),
				),
			},
			want: types.NewUnion(
				types.NewNamed("InlineKeyboardMarkup", types.KindObject),
				types.NewNamed("ReplyKeyboardMarkup", types.KindObject),
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
			got, err := overlays.ReplyMarkup{}.Apply(tc.field).Type()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ReplyMarkup.Apply must return the original field unchanged when its Type method errors",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ReplyMarkup.Apply must replace the four-variant markup union with ReplyMarkup and pass through all other types",
			)
		})
	}
}
