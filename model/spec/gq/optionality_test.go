// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec/gq"
)

func TestObjectFieldOptionality_Value(t *testing.T) {
	cases := []struct {
		name    string
		text    string
		want    model.Optionality
		wantErr bool
	}{
		{
			name: "returns true when description starts with Optional. prefix",
			text: "Optional. The message was edited.",
			want: true,
		},
		{
			name: "returns false when description has no Optional. prefix",
			text: "Unique identifier of the target chat.",
			want: false,
		},
		{
			name: "returns false when description starts with Optional. without trailing space",
			text: "Optional.Something",
			want: false,
		},
		{
			name:    "returns error when the selection is empty",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := emptySel()
			if !tc.wantErr {
				sel = tdWith(tc.text)
			}
			got, err := gq.NewObjectFieldOptionality(sel).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ObjectFieldOptionality.Value must return an error when the td selection is empty",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ObjectFieldOptionality.Value must return true only when the description starts with the exact prefix 'Optional. '",
			)
		})
	}
}

func TestMethodFieldOptionality_Value(t *testing.T) {
	cases := []struct {
		name    string
		text    string
		want    model.Optionality
		wantErr bool
	}{
		{
			name: "returns true when the required column contains exactly Optional",
			text: "Optional",
			want: true,
		},
		{
			name: "returns false when the required column contains Yes",
			text: "Yes",
			want: false,
		},
		{
			name: "returns false when the text is Optional with trailing content",
			text: "Optional. some note",
			want: false,
		},
		{
			name:    "returns error when the selection is empty",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := emptySel()
			if !tc.wantErr {
				sel = tdWith(tc.text)
			}
			got, err := gq.NewMethodFieldOptionality(sel).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"MethodFieldOptionality.Value must return an error when the td selection is empty",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"MethodFieldOptionality.Value must return true only when the required column contains exactly 'Optional'",
			)
		})
	}
}
