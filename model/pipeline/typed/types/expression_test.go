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
	"github.com/andreychh/tgen/model/typeexpr"
)

func TestExpression_Value(t *testing.T) {
	cases := []struct {
		name    string
		phrase  prose.Phrase
		want    typeexpr.Expression
		wantErr string
	}{
		{
			name:   "decodes the two built-in types a chat identifier is written as",
			phrase: prose.NewPhrase(plain("Integer or String")),
			want: typeexpr.NewUnion(
				typeexpr.NewPrimitive(primitive.Integer),
				typeexpr.NewPrimitive(primitive.String),
			),
		},
		{
			name:   "decodes a sequence of a documented type",
			phrase: prose.NewPhrase(plain("Array of "), anchor("MessageEntity")),
			want:   typeexpr.NewArray(typeexpr.NewNamed(model.Reference("messageentity"))),
		},
		{
			name:   "decodes the nested sequence a keyboard is written as",
			phrase: prose.NewPhrase(plain("Array of Array of "), anchor("InlineKeyboardButton")),
			want: typeexpr.NewArray(
				typeexpr.NewArray(typeexpr.NewNamed(model.Reference("inlinekeyboardbutton"))),
			),
		},
		{
			name: "decodes a sequence whose element the documentation lists out",
			phrase: prose.NewPhrase(
				plain("Array of "),
				anchor("InputMediaAudio"),
				plain(", "),
				anchor("InputMediaPhoto"),
				plain(" and "),
				anchor("InputMediaVideo"),
			),
			want: typeexpr.NewArray(typeexpr.NewUnion(
				typeexpr.NewNamed(model.Reference("inputmediaaudio")),
				typeexpr.NewNamed(model.Reference("inputmediaphoto")),
				typeexpr.NewNamed(model.Reference("inputmediavideo")),
			)),
		},
		{
			name:   "decodes a choice between a documented type and a built-in one",
			phrase: prose.NewPhrase(anchor("InputFile"), plain(" or String")),
			want: typeexpr.NewUnion(
				typeexpr.NewNamed(model.Reference("inputfile")),
				typeexpr.NewPrimitive(primitive.String),
			),
		},
		{
			name:    "names the lexing stage when the prose carries a word it cannot read",
			phrase:  prose.NewPhrase(plain("Целое")),
			wantErr: "lexing type expression",
		},
		{
			name:    "names the parsing stage when the tokens spell no single type",
			phrase:  prose.NewPhrase(plain("Integer String")),
			wantErr: "parsing type expression",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := types.NewExpression(tc.phrase).Value()
			if tc.wantErr != "" {
				assert.ErrorContains(
					t,
					err,
					tc.wantErr,
					"Expression.Value must name the stage that refused the prose",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Expression.Value must decode a type cell into the type expression its prose spells",
			)
		})
	}
}
