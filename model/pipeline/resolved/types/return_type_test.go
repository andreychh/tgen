// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/resolved/types"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

func TestReturnType_Value(t *testing.T) {
	message := typeexpr.NewNamed(model.Reference("message"))
	cases := []struct {
		name        string
		description prose.Passage
		want        typeexpr.Expression
		wantErr     bool
	}{
		{
			name: "returns the documented type a Returns clause names",
			description: passage(
				plain("Returns "),
				anchor("BotName"),
				plain(" on success."),
			),
			want: typeexpr.NewNamed(model.Reference("botname")),
		},
		{
			name: "returns the built-in type a Returns clause names",
			description: passage(
				plain("Returns the created invoice link as "),
				italic("String"),
				plain(" on success."),
			),
			want: typeexpr.NewPrimitive(primitive.String),
		},
		{
			name: "returns the documented type an is-returned clause names",
			description: passage(
				plain("On success, the sent "),
				anchor("Message"),
				plain(" is returned."),
			),
			want: message,
		},
		{
			name: "returns the built-in type an is-returned clause names",
			description: passage(
				italic("True"),
				plain(" is returned on success."),
			),
			want: typeexpr.NewPrimitive(primitive.True),
		},
		{
			name: "wraps the documented type an Array of clause names in an array",
			description: passage(
				plain("Returns an Array of "),
				anchor("Sticker"),
				plain(" objects."),
			),
			want: typeexpr.NewArray(typeexpr.NewNamed(model.Reference("sticker"))),
		},
		{
			name: "reads an Array of clause as a sequence and not as the type the Returns word alone would name",
			description: passage(
				plain("Returns an Array of "),
				anchor("BotCommand"),
				plain(" objects. If commands aren't set, an empty list is returned."),
			),
			want: typeexpr.NewArray(typeexpr.NewNamed(model.Reference("botcommand"))),
		},
		{
			name: "joins both branches of a conditional clause rather than stopping at the first",
			description: passage(
				plain("On success, if the edited message is not an inline message, the edited "),
				anchor("Message"),
				plain(" is returned, otherwise "),
				italic("True"),
				plain(" is returned."),
			),
			want: typeexpr.NewUnion(message, typeexpr.NewPrimitive(primitive.True)),
		},
		{
			name: "reads only the first branch when the second is not introduced by otherwise",
			description: passage(
				anchor("Message"),
				plain(" is returned, or else "),
				italic("True"),
				plain(" is returned."),
			),
			want: message,
		},
		{
			name: "reads only the first branch when the second names a documented type rather than a built-in one",
			description: passage(
				anchor("Message"),
				plain(" is returned, otherwise "),
				anchor("Story"),
				plain(" is returned."),
			),
			want: message,
		},
		{
			name: "finds a clause that opens after the words leading up to it",
			description: passage(
				plain("Use this method to send text messages. On success, the sent "),
				anchor("Message"),
				plain(" is returned."),
			),
			want: message,
		},
		{
			name: "finds a clause that opens after a line break",
			description: passage(
				plain("On success:"),
				prose.NewLineBreak(),
				plain("Returns "),
				anchor("Story"),
				plain("."),
			),
			want: typeexpr.NewNamed(model.Reference("story")),
		},
		{
			name: "ignores the space around a built-in type word",
			description: passage(
				plain("Returns "),
				italic(" Integer "),
				plain(" on success."),
			),
			want: typeexpr.NewPrimitive(primitive.Integer),
		},
		{
			name: "reads the clause the first paragraph carries, leaving the later ones alone",
			description: prose.NewPassage(
				prose.NewParagraph(plain("Returns "), anchor("Story"), plain(".")),
				prose.NewParagraph(plain("Returns "), anchor("Message"), plain(".")),
			),
			want: typeexpr.NewNamed(model.Reference("story")),
		},
		{
			name: "returns error when the emphasized word names no built-in type",
			description: passage(
				plain("Returns "),
				italic("Widget"),
				plain(" on success."),
			),
			wantErr: true,
		},
		{
			name: "returns error when an Array of clause names no documented type",
			description: passage(
				plain("An Array of "),
				italic("Widget"),
				plain(" is sent."),
			),
			wantErr: true,
		},
		{
			name: "returns error when the link leaves the documentation page",
			description: passage(
				plain("Returns "),
				prose.NewLink("Message", prose.StylePlain, "https://core.telegram.org/bots/api"),
				plain(" on success."),
			),
			wantErr: true,
		},
		{
			name: "returns error when the paragraph names no type at all",
			description: passage(
				plain("Use this method to log out from the cloud Bot API server."),
			),
			wantErr: true,
		},
		{
			name: "returns error when the clause closes the paragraph with its type missing",
			description: passage(
				plain("On success, nothing "),
				plain("is returned."),
			),
			wantErr: true,
		},
		{
			name: "returns error when the description opens with a list",
			description: prose.NewPassage(
				prose.NewList(prose.NewItem(plain("Returns "), anchor("Story"))),
			),
			wantErr: true,
		},
		{
			name:        "returns error when the description carries no blocks",
			description: prose.NewPassage(),
			wantErr:     true,
		},
		{
			name:        "returns error when the first paragraph carries no inlines",
			description: passage(),
			wantErr:     true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := types.NewReturnType(tc.description).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ReturnType.Value must refuse a description whose first paragraph spells no return type",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ReturnType.Value must decode the return clause into the type it names, the array and union forms before the plain ones",
			)
		})
	}
}
