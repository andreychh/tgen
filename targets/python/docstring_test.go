// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tpython "github.com/andreychh/tgen/targets/python"
)

func TestDocString_Value(t *testing.T) {
	cases := []struct {
		name   string
		docstr tpython.DocString
		want   string
	}{
		{
			name:   "wraps short content inline with triple-quote delimiters",
			docstr: tpython.NewClassDocString("Hello world."),
			want:   `"""Hello world."""`,
		},
		{
			name: "wraps content of exactly 72 characters inline without multiline format",
			docstr: tpython.NewClassDocString(
				"123456789012345678901234567890123456789012345678901234567890123456789012",
			),
			want: `"""123456789012345678901234567890123456789012345678901234567890123456789012"""`,
		},
		{
			name: "switches to multiline format and word-wraps when content exceeds 72 characters",
			docstr: tpython.NewClassDocString(
				"This description is long enough to require word wrapping because it exceeds the seventy-two character limit.",
			),
			want: "\"\"\"\n" +
				"    This description is long enough to require word wrapping because it\n" +
				"    exceeds the seventy-two character limit.\n" +
				"    \"\"\"",
		},
		{
			name:   "places short two-paragraph content inline preserving the paragraph break",
			docstr: tpython.NewClassDocString("First paragraph.\n\nSecond paragraph."),
			want:   "\"\"\"First paragraph.\n\nSecond paragraph.\"\"\"",
		},
		{
			name: "separates paragraphs with a blank line in multiline format",
			docstr: tpython.NewClassDocString(
				"This is the first paragraph that is definitely longer than seventy-two characters in total.\n\nSecond paragraph.",
			),
			want: "\"\"\"\n" +
				"    This is the first paragraph that is definitely longer than seventy-two\n" +
				"    characters in total.\n\n" +
				"    Second paragraph.\n" +
				"    \"\"\"",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.docstr.Value(),
				"DocString.Value must use inline format for content up to 72 characters and multiline with word-wrap and indentation for longer content",
			)
		})
	}
}
