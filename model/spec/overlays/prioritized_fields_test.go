// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/spec/overlays"
)

func collectKeys(seq func(yield func(spec.Field) bool)) []model.Key {
	var keys []model.Key
	for f := range seq {
		k, _ := f.Key()
		keys = append(keys, k)
	}
	return keys
}

func TestNewPrioritizedFields(t *testing.T) {
	req := func(key model.Key) spec.Field { return stubField{key: key, optional: false} }
	opt := func(key model.Key) spec.Field { return stubField{key: key, optional: true} }
	errF := func(key model.Key) spec.Field { return stubField{key: key, optErr: errOpt} }
	cases := []struct {
		name   string
		fields []spec.Field
		want   []model.Key
	}{
		{
			name:   "emits nothing when the input sequence is empty",
			fields: nil,
			want:   nil,
		},
		{
			name:   "preserves original order when all fields are required",
			fields: []spec.Field{req("a"), req("b"), req("c")},
			want:   []model.Key{"a", "b", "c"},
		},
		{
			name:   "preserves original order when all fields are optional",
			fields: []spec.Field{opt("x"), opt("y")},
			want:   []model.Key{"x", "y"},
		},
		{
			name:   "emits required fields before optional fields regardless of original order",
			fields: []spec.Field{opt("first"), req("second"), opt("third"), req("fourth")},
			want:   []model.Key{"second", "fourth", "first", "third"},
		},
		{
			name:   "treats fields with erroring Optionality as required and emits them first",
			fields: []spec.Field{opt("optional"), errF("unknown")},
			want:   []model.Key{"unknown", "optional"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := collectKeys(overlays.NewPrioritizedFields(slices.Values(tc.fields)))
			assert.Equal(
				t,
				tc.want,
				got,
				"NewPrioritizedFields must emit all required fields before all optional fields",
			)
		})
	}
}
