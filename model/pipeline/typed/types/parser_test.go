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
	"github.com/andreychh/tgen/model/typeexpr"
)

func TestParser_Expression(t *testing.T) {
	audio := types.NewRef(model.Reference("inputmediaaudio"))
	photo := types.NewRef(model.Reference("inputmediaphoto"))
	video := types.NewRef(model.Reference("inputmediavideo"))
	namedAudio := typeexpr.NewNamed(model.Reference("inputmediaaudio"))
	namedPhoto := typeexpr.NewNamed(model.Reference("inputmediaphoto"))
	namedVideo := typeexpr.NewNamed(model.Reference("inputmediavideo"))
	integer := types.NewPrimitive(primitive.Integer)
	str := types.NewPrimitive(primitive.String)
	cases := []struct {
		name    string
		tokens  []types.Token
		want    typeexpr.Expression
		wantErr bool
	}{
		{
			name:   "returns the documented type a lone reference names",
			tokens: []types.Token{audio},
			want:   namedAudio,
		},
		{
			name:   "returns the built-in type a lone primitive names",
			tokens: []types.Token{integer},
			want:   typeexpr.NewPrimitive(primitive.Integer),
		},
		{
			name:   "joins the terms an or separates into one union",
			tokens: []types.Token{integer, types.NewOr(), str},
			want: typeexpr.NewUnion(
				typeexpr.NewPrimitive(primitive.Integer),
				typeexpr.NewPrimitive(primitive.String),
			),
		},
		{
			name:   "flattens a union of more than two terms into one level",
			tokens: []types.Token{audio, types.NewOr(), photo, types.NewOr(), video},
			want:   typeexpr.NewUnion(namedAudio, namedPhoto, namedVideo),
		},
		{
			name:   "wraps the term an Array of precedes in an array",
			tokens: []types.Token{types.NewArrayOf(), photo},
			want:   typeexpr.NewArray(namedPhoto),
		},
		{
			name:   "nests an array inside an array for a repeated Array of",
			tokens: []types.Token{types.NewArrayOf(), types.NewArrayOf(), photo},
			want:   typeexpr.NewArray(typeexpr.NewArray(namedPhoto)),
		},
		{
			name: "keeps a list an array introduces inside that array",
			tokens: []types.Token{
				types.NewArrayOf(),
				audio,
				types.NewSeries(),
				photo,
				types.NewSeries(),
				video,
			},
			want: typeexpr.NewArray(typeexpr.NewUnion(namedAudio, namedPhoto, namedVideo)),
		},
		{
			name:   "lifts an or out of the array that precedes it",
			tokens: []types.Token{types.NewArrayOf(), photo, types.NewOr(), video},
			want:   typeexpr.NewUnion(typeexpr.NewArray(namedPhoto), namedVideo),
		},
		{
			name:    "returns error when a list separator stands outside an array",
			tokens:  []types.Token{audio, types.NewSeries(), photo},
			wantErr: true,
		},
		{
			name:    "returns error when a separator opens the stream",
			tokens:  []types.Token{types.NewOr(), audio},
			wantErr: true,
		},
		{
			name:    "returns error when a separator closes the stream",
			tokens:  []types.Token{audio, types.NewOr()},
			wantErr: true,
		},
		{
			name:    "returns error when two terms stand with no separator between them",
			tokens:  []types.Token{audio, photo},
			wantErr: true,
		},
		{
			name:    "returns error when an Array of closes the stream",
			tokens:  []types.Token{types.NewArrayOf()},
			wantErr: true,
		},
		{
			name:    "returns error when the stream carries no tokens",
			tokens:  nil,
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := types.NewParser(tc.tokens).Expression()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Parser.Expression must refuse a token stream that spells no single type",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Parser.Expression must bind an or across whole terms and a list inside the array it opens",
			)
		})
	}
}

func TestParser_ExpressionRepeated(t *testing.T) {
	parser := types.NewParser([]types.Token{
		types.NewPrimitive(primitive.Integer),
		types.NewOr(),
		types.NewPrimitive(primitive.String),
	})
	first, err := parser.Expression()
	require.NoError(t, err)
	second, err := parser.Expression()
	require.NoError(t, err)
	assert.Equal(
		t,
		first,
		second,
		"Parser.Expression must decode the same stream on every call, holding no position between them",
	)
}
