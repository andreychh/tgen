// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/typed/types"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
)

func TestLexer_Tokens(t *testing.T) {
	cases := []struct {
		name    string
		phrase  prose.Phrase
		want    []types.Token
		wantErr bool
	}{
		{
			name:   "lexes a built-in type word into a primitive token",
			phrase: prose.NewPhrase(plain("Integer")),
			want:   []types.Token{types.NewPrimitive(primitive.Integer)},
		},
		{
			name:   "lexes an alias of a built-in type word into the kind it spells",
			phrase: prose.NewPhrase(plain("Int")),
			want:   []types.Token{types.NewPrimitive(primitive.Integer)},
		},
		{
			name:   "lexes an anchor link into a reference token",
			phrase: prose.NewPhrase(anchor("PhotoSize")),
			want:   []types.Token{types.NewRef(model.Reference("photosize"))},
		},
		{
			name:   "lexes the two words of the Array of prefix into one token",
			phrase: prose.NewPhrase(plain("Array of "), anchor("PhotoSize")),
			want: []types.Token{
				types.NewArrayOf(),
				types.NewRef(model.Reference("photosize")),
			},
		},
		{
			name:   "lexes or into the separator that joins whole terms",
			phrase: prose.NewPhrase(plain("Integer or String")),
			want: []types.Token{
				types.NewPrimitive(primitive.Integer),
				types.NewOr(),
				types.NewPrimitive(primitive.String),
			},
		},
		{
			name: "lexes a comma and a closing and into list separators",
			phrase: prose.NewPhrase(
				anchor("InputMediaAudio"),
				plain(", "),
				anchor("InputMediaPhoto"),
				plain(" and "),
				anchor("InputMediaVideo"),
			),
			want: []types.Token{
				types.NewRef(model.Reference("inputmediaaudio")),
				types.NewSeries(),
				types.NewRef(model.Reference("inputmediaphoto")),
				types.NewSeries(),
				types.NewRef(model.Reference("inputmediavideo")),
			},
		},
		{
			name:   "lexes a comma that no space surrounds",
			phrase: prose.NewPhrase(plain("Integer,String")),
			want: []types.Token{
				types.NewPrimitive(primitive.Integer),
				types.NewSeries(),
				types.NewPrimitive(primitive.String),
			},
		},
		{
			name:   "lexes nothing from a phrase carrying no inlines",
			phrase: prose.NewPhrase(),
			want:   nil,
		},
		{
			name: "returns error when a link addresses somewhere other than the page",
			phrase: prose.NewPhrase(
				prose.NewLink("File", prose.StylePlain, "https://example.org/file"),
			),
			wantErr: true,
		},
		{
			name:    "returns error when Array is followed by a word other than of",
			phrase:  prose.NewPhrase(plain("Array in String")),
			wantErr: true,
		},
		{
			name:    "returns error when Array closes the text run",
			phrase:  prose.NewPhrase(plain("Array")),
			wantErr: true,
		},
		{
			name:    "returns error when a word names no built-in type",
			phrase:  prose.NewPhrase(plain("Строка")),
			wantErr: true,
		},
		{
			name:    "returns error when an inline is neither text nor a link",
			phrase:  prose.NewPhrase(prose.NewLineBreak()),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := types.NewLexer(tc.phrase).Tokens()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Lexer.Tokens must refuse prose it cannot turn into type tokens",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Lexer.Tokens must yield one token per lexical unit, in the order the prose writes them",
			)
		})
	}
}
