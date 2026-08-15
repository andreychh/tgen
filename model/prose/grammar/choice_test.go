// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package grammar_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/prose/grammar"
)

func TestChoice_Match(t *testing.T) {
	cases := []struct {
		name    string
		rules   []grammar.Rule[string]
		inlines []prose.Inline
		want    string
		wantOK  bool
	}{
		{
			name:    "returns what the only alternative decodes",
			rules:   []grammar.Rule[string]{newWord("Returns", "returns")},
			inlines: []prose.Inline{plain("Returns")},
			want:    "returns",
			wantOK:  true,
		},
		{
			name: "returns what a later alternative decodes when the ones before it decline",
			rules: []grammar.Rule[string]{
				newWord("Array of", "array"),
				newWord("otherwise", "union"),
				newWord("Returns", "returns"),
			},
			inlines: []prose.Inline{plain("Returns")},
			want:    "returns",
			wantOK:  true,
		},
		{
			name: "returns what the first matching alternative decodes, leaving the rest untried",
			rules: []grammar.Rule[string]{
				newWord("Returns", "array"),
				newWord("Returns", "returns"),
			},
			inlines: []prose.Inline{plain("Returns")},
			want:    "array",
			wantOK:  true,
		},
		{
			name: "returns nothing when every alternative declines",
			rules: []grammar.Rule[string]{
				newWord("Array of", "array"),
				newWord("otherwise", "union"),
			},
			inlines: []prose.Inline{plain("Returns")},
			want:    "",
			wantOK:  false,
		},
		{
			name:    "returns nothing when the choice holds no alternatives",
			rules:   nil,
			inlines: []prose.Inline{plain("Returns")},
			want:    "",
			wantOK:  false,
		},
		{
			name:    "returns nothing when the run holds no inlines",
			rules:   []grammar.Rule[string]{newWord("Returns", "returns")},
			inlines: nil,
			want:    "",
			wantOK:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := grammar.NewChoice(tc.rules...).Match(tc.inlines)
			if !tc.wantOK {
				assert.False(
					t,
					ok,
					"Choice.Match must report no match when no alternative recognizes the run",
				)
				return
			}
			require.True(t, ok)
			assert.Equal(
				t,
				tc.want,
				got,
				"Choice.Match must return what the earliest matching alternative decodes",
			)
		})
	}
}
