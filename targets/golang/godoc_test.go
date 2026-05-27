// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tgolang "github.com/andreychh/tgen/targets/golang"
)

func TestGoDoc_Value(t *testing.T) {
	cases := []struct {
		name  string
		godoc tgolang.GoDoc
		want  string
	}{
		{
			name:  "prefixes a single short line with the comment marker",
			godoc: tgolang.NewTypeGodoc("Hello world."),
			want:  "// Hello world.",
		},
		{
			name: "wraps content at 77 characters before adding the comment prefix",
			godoc: tgolang.NewTypeGodoc(
				"This is a very long description that will exceed the seventy-seven character limit and must be wrapped.",
			),
			want: "// This is a very long description that will exceed the seventy-seven character\n// limit and must be wrapped.",
		},
		{
			name:  "separates two paragraphs with an empty comment line",
			godoc: tgolang.NewTypeGodoc("First paragraph.\n\nSecond paragraph."),
			want:  "// First paragraph.\n//\n// Second paragraph.",
		},
		{
			name: "indents continuation lines within a paragraph for a field-level comment",
			godoc: tgolang.NewFieldGodoc(
				"This description is intentionally long so that the word wrapper splits it across two lines for testing purposes.",
			),
			want: "// This description is intentionally long so that the word wrapper splits it\n\t// across two lines for testing purposes.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				tc.godoc.Value(),
				"GoDoc.Value must wrap lines at 77 characters and prefix every line with '// ', separating paragraphs with an empty comment line",
			)
		})
	}
}
