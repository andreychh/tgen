// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/spec/gq"
	"github.com/andreychh/tgen/model/types"
)

func TestMethod_Result(t *testing.T) {
	msg := types.NewNamed("Message", types.KindObject)
	trueType := types.NewNamed("True", types.KindPrimitive)
	cases := []struct {
		name    string
		desc    string
		noDesc  bool
		want    spec.Result
		wantErr bool
	}{
		{
			name: "classifies a bare True success sentinel as a Command",
			desc: "Returns True on success.",
			want: spec.NewCommand(),
		},
		{
			name: "classifies a named return as a Value carrying the expression",
			desc: "Returns Message on success.",
			want: spec.NewValue(msg),
		},
		{
			name: "classifies a union containing True as a Value, not a Command",
			desc: "the Message is returned, otherwise True is returned",
			want: spec.NewValue(types.NewUnion(msg, trueType)),
		},
		{
			name:    "propagates the error when the return type cannot be extracted",
			noDesc:  true,
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, h4 := returnTypeFixture(tc.desc)
			if tc.noDesc {
				root, h4 = returnTypeNoParaFixture()
			}
			got, err := gq.NewMethod(root, h4).Result()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Method.Result must propagate errors from return type extraction",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Method.Result must classify a bare True as a Command and any other return as a Value",
			)
		})
	}
}
