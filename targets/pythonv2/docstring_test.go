// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/targets/pythonv2"
)

func plain(content string) prose.Text {
	return prose.NewText(content, prose.StylePlain)
}

func TestDocstring_Value(t *testing.T) {
	cases := []struct {
		name    string
		passage prose.Passage
		want    string
	}{
		{
			name:    "writes nothing for a passage holding no block",
			passage: prose.NewPassage(),
			want:    "",
		},
		{
			name:    "writes nothing for a paragraph reading as nothing",
			passage: prose.NewPassage(prose.NewParagraph(plain("   "))),
			want:    "",
		},
		{
			name:    "closes a short paragraph on the line it opened",
			passage: prose.NewPassage(prose.NewParagraph(plain("Unique identifier."))),
			want:    `"""Unique identifier."""`,
		},
		{
			name:    "keeps a non-ASCII paragraph whole",
			passage: prose.NewPassage(prose.NewParagraph(plain("Идентификатор чата — «канал»."))),
			want:    `"""Идентификатор чата — «канал»."""`,
		},
		{
			name: "reads a link by its text and not by what it addresses",
			passage: prose.NewPassage(
				prose.NewParagraph(
					plain("Sender, see "),
					prose.NewLink("User", prose.StylePlain, "#user"),
				),
			),
			want: `"""Sender, see User"""`,
		},
		{
			name:    "closes on a line of its own where the prose ends in a quote",
			passage: prose.NewPassage(prose.NewParagraph(plain(`Currently one of "photo"`))),
			want: `"""Currently one of "photo"
    """`,
		},
		{
			name: "folds a paragraph outgrowing the room a line has, indenting every line after the first",
			passage: prose.NewPassage(prose.NewParagraph(plain(
				"Unique identifier for the target chat or username of the target channel in the format channelusername",
			))),
			want: `"""Unique identifier for the target chat or username of the target
    channel in the format channelusername
    """`,
		},
		{
			name: "parts two paragraphs with a line holding nothing at all",
			passage: prose.NewPassage(
				prose.NewParagraph(plain("Describes the first thing.")),
				prose.NewParagraph(plain("Describes the second thing.")),
			),
			want: `"""Describes the first thing.

    Describes the second thing.
    """`,
		},
		{
			name: "writes a list as one dashed line per item",
			passage: prose.NewPassage(
				prose.NewParagraph(plain("One of:")),
				prose.NewList(
					prose.NewItem(plain("a file identifier")),
					prose.NewItem(plain("an HTTP URL")),
				),
			),
			want: `"""One of:

    - a file identifier
    - an HTTP URL
    """`,
		},
		{
			name: "starts a line of its own at a forced break",
			passage: prose.NewPassage(prose.NewParagraph(
				plain("First line."),
				prose.NewLineBreak(),
				plain("Second line."),
			)),
			want: `"""First line.
    Second line.
    """`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := pythonv2.NewClassDocstring(c.passage).Value()
			assert.Equal(
				t,
				c.want,
				got,
				"Docstring must enclose the prose in triple quotes and indent every line after the first",
			)
		})
	}
}

func TestDocstring_ValueForField(t *testing.T) {
	phrase := prose.NewPhrase(plain("Optional. Caption of the document."))
	got := pythonv2.NewFieldDocstring(phrase).Value()
	assert.Equal(
		t,
		`"""Optional. Caption of the document."""`,
		got,
		"Docstring must read the one phrase describing a field as a single paragraph",
	)
}
