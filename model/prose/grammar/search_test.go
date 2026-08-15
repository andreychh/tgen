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

func TestSearch_Find(t *testing.T) {
	cases := []struct {
		name    string
		rule    grammar.Rule[string]
		inlines []prose.Inline
		want    string
		wantOK  bool
	}{
		{
			name:    "returns what the rule decodes at the front of the run",
			rule:    newWord("Returns", "returns"),
			inlines: []prose.Inline{plain("Returns"), anchor("Message")},
			want:    "returns",
			wantOK:  true,
		},
		{
			name: "returns what the rule decodes further along, skipping the positions before it",
			rule: newWord("Returns", "returns"),
			inlines: []prose.Inline{
				plain("On success"),
				anchor("Story"),
				plain("Returns"),
			},
			want:   "returns",
			wantOK: true,
		},
		{
			name: "returns what the rule decodes at the earliest of two matching positions",
			rule: newWord("Returns", "returns"),
			inlines: []prose.Inline{
				plain("Returns"),
				plain("Returns"),
			},
			want:   "returns",
			wantOK: true,
		},
		{
			name: "returns what the rule decodes at the last position of the run",
			rule: newWord("is returned", "returned"),
			inlines: []prose.Inline{
				plain("On success"),
				plain("is returned"),
			},
			want:   "returned",
			wantOK: true,
		},
		{
			name:    "returns nothing when the rule matches at no position",
			rule:    newWord("Returns", "returns"),
			inlines: []prose.Inline{plain("On success"), anchor("Message")},
			want:    "",
			wantOK:  false,
		},
		{
			name:    "returns nothing when the run holds no inlines",
			rule:    newWord("Returns", "returns"),
			inlines: nil,
			want:    "",
			wantOK:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := grammar.NewSearch(tc.rule).Find(tc.inlines)
			if !tc.wantOK {
				assert.False(
					t,
					ok,
					"Search.Find must report no match when the rule recognizes no position of the run",
				)
				return
			}
			require.True(t, ok)
			assert.Equal(
				t,
				tc.want,
				got,
				"Search.Find must return what the rule decodes at the earliest position it matches",
			)
		})
	}
}
