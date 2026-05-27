// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/spec/overlays"
	"github.com/andreychh/tgen/model/types"
)

func TestMaybeMessage_ReturnType(t *testing.T) {
	message := types.NewNamed("Message", types.KindObject)
	trueType := types.NewNamed("True", types.KindPrimitive)
	cases := []struct {
		name    string
		method  spec.Method
		want    types.Expression
		wantErr bool
	}{
		{
			name:   "replaces Message or True union return type with MaybeMessage named type",
			method: stubMethod{returnType: types.NewUnion(message, trueType)},
			want:   types.NewNamed("MaybeMessage", types.KindUnion),
		},
		{
			name:   "passes through Message return type without replacement",
			method: stubMethod{returnType: message},
			want:   message,
		},
		{
			name: "passes through a non-matching union return type without replacement",
			method: stubMethod{
				returnType: types.NewUnion(message, types.NewNamed("Chat", types.KindObject)),
			},
			want: types.NewUnion(message, types.NewNamed("Chat", types.KindObject)),
		},
		{
			name:    "propagates error when the inner method's ReturnType returns an error",
			method:  stubMethod{returnErr: errors.New("no return type")},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlays.NewMaybeMessage(tc.method).ReturnType()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"MaybeMessage.ReturnType must propagate errors from the inner method",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"MaybeMessage.ReturnType must replace Message|True with MaybeMessage and pass through all other return types",
			)
		})
	}
}
