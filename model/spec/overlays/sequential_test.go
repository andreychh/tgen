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

func TestSequential_Apply(t *testing.T) {
	intType := types.NewNamed("Integer", types.KindPrimitive)
	strType := types.NewNamed("String", types.KindPrimitive)
	chatID := types.NewNamed("ChatID", types.KindUnion)
	cases := []struct {
		name  string
		items []overlays.Overlay
		field spec.Field
		want  types.Expression
	}{
		{
			name:  "passes field through unchanged when no overlays are registered",
			items: nil,
			field: stubField{expr: intType},
			want:  intType,
		},
		{
			name:  "applies a single overlay to the field",
			items: []overlays.Overlay{overlays.ChatID{}},
			field: stubField{expr: types.NewUnion(intType, strType)},
			want:  chatID,
		},
		{
			name: "applies overlays in registration order so each receives the previous result",
			items: []overlays.Overlay{
				overlays.ChatID{},
				overlays.ReplyMarkup{},
			},
			field: stubField{expr: types.NewUnion(intType, strType)},
			want:  chatID,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seq := overlays.NewSequential(tc.items...)
			got, err := seq.Apply(tc.field).Type()
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Sequential.Apply must delegate to each overlay in order and return the final field",
			)
		})
	}
}
