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

func TestMaybeMessage_Result(t *testing.T) {
	message := types.NewNamed("Message", types.KindObject)
	trueType := types.NewNamed("True", types.KindPrimitive)
	chat := types.NewNamed("Chat", types.KindObject)
	cases := []struct {
		name    string
		method  spec.Method
		want    spec.Result
		wantErr bool
	}{
		{
			name:   "replaces a Message or True Value with a MaybeMessage Value",
			method: stubMethod{result: spec.NewValue(types.NewUnion(message, trueType))},
			want:   spec.NewValue(types.NewNamed("MaybeMessage", types.KindUnion)),
		},
		{
			name:   "passes through a Message Value without replacement",
			method: stubMethod{result: spec.NewValue(message)},
			want:   spec.NewValue(message),
		},
		{
			name: "passes through a non-matching union Value without replacement",
			method: stubMethod{
				result: spec.NewValue(types.NewUnion(message, chat)),
			},
			want: spec.NewValue(types.NewUnion(message, chat)),
		},
		{
			name:   "passes through a Command without replacement",
			method: stubMethod{result: spec.NewCommand()},
			want:   spec.NewCommand(),
		},
		{
			name:    "propagates error when the inner method's Result returns an error",
			method:  stubMethod{resultErr: errors.New("no result")},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlays.NewMaybeMessage(tc.method).Result()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"MaybeMessage.Result must propagate errors from the inner method",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"MaybeMessage.Result must replace a Message|True Value with a MaybeMessage Value and pass through all other results",
			)
		})
	}
}
