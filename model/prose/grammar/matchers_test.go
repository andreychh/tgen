// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package grammar_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/prose/grammar"
)

func TestMarker_Matches(t *testing.T) {
	returns := regexp.MustCompile(`(?i)\breturns\b`)
	trailing := regexp.MustCompile(`(?i)\bmust be\s*$`)
	cases := []struct {
		name    string
		pattern *regexp.Regexp
		inline  prose.Inline
		want    bool
	}{
		{
			name:    "reports true when a plain run carries the signal word",
			pattern: returns,
			inline:  plain("Returns the sent message on success."),
			want:    true,
		},
		{
			name:    "reports true whatever case the signal word is written in",
			pattern: returns,
			inline:  plain("RETURNS the sent message."),
			want:    true,
		},
		{
			name:    "reports true when the signal closes the run the pattern anchors to",
			pattern: trailing,
			inline:  plain("Type of the result, must be "),
			want:    true,
		},
		{
			name:    "reports false when the anchored signal does not close the run",
			pattern: trailing,
			inline:  plain("must be one of the values below"),
			want:    false,
		},
		{
			name:    "reports false when the run carries the signal but is set in italic",
			pattern: returns,
			inline:  italic("Returns"),
			want:    false,
		},
		{
			name:    "reports false when the run carries the signal but is set in bold",
			pattern: returns,
			inline:  prose.NewText("Returns", prose.StyleBold),
			want:    false,
		},
		{
			name:    "reports false when a link carries the signal in its text",
			pattern: returns,
			inline:  prose.NewLink("Returns", prose.StylePlain, "#message"),
			want:    false,
		},
		{
			name:    "reports false when the inline is a line break",
			pattern: returns,
			inline:  prose.NewLineBreak(),
			want:    false,
		},
		{
			name:    "reports false when the signal word is only part of a longer word",
			pattern: returns,
			inline:  plain("The value returnsheet names."),
			want:    false,
		},
		{
			name:    "reports false when the run is empty",
			pattern: returns,
			inline:  plain(""),
			want:    false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := grammar.NewMarker(tc.pattern).Matches(tc.inline)
			assert.Equal(
				t,
				tc.want,
				got,
				"Marker.Matches must accept a plain run its pattern satisfies and nothing else",
			)
		})
	}
}

func TestCapture_Matches(t *testing.T) {
	always := regexp.MustCompile(`(?i)\balways\s*“([^”]+)”\s*$`)
	pair := regexp.MustCompile(`(\w+)-(\w+)`)
	groupless := regexp.MustCompile(`(?i)\balways\s*“[^”]+”\s*$`)
	cases := []struct {
		name    string
		pattern *regexp.Regexp
		inline  prose.Inline
		want    string
		wantOK  bool
	}{
		{
			name:    "returns the value the first group lifts out of a plain run",
			pattern: always,
			inline:  plain("The member's status in the chat, always “creator”"),
			want:    "creator",
			wantOK:  true,
		},
		{
			name:    "returns a value written outside the ASCII range",
			pattern: always,
			inline:  plain("Тип результата, always “видео_заметка”"),
			want:    "видео_заметка",
			wantOK:  true,
		},
		{
			name:    "returns the first group when the pattern holds more than one",
			pattern: pair,
			inline:  plain("kind: hidden-user"),
			want:    "hidden",
			wantOK:  true,
		},
		{
			name:    "returns nothing when the pattern captures no group at all",
			pattern: groupless,
			inline:  plain("The member's status in the chat, always “creator”"),
			want:    "",
			wantOK:  false,
		},
		{
			name:    "returns nothing when the quoted value does not close the run",
			pattern: always,
			inline:  plain("always “creator” of the chat"),
			want:    "",
			wantOK:  false,
		},
		{
			name:    "returns nothing when the run is set in italic",
			pattern: always,
			inline:  italic("always “creator”"),
			want:    "",
			wantOK:  false,
		},
		{
			name:    "returns nothing when a link carries the quoted value",
			pattern: always,
			inline:  prose.NewLink("always “creator”", prose.StylePlain, "#chatmemberowner"),
			want:    "",
			wantOK:  false,
		},
		{
			name:    "returns nothing when the pattern does not match the run",
			pattern: always,
			inline:  plain("Unique identifier of the message"),
			want:    "",
			wantOK:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := grammar.NewCapture(tc.pattern).Matches(tc.inline)
			if !tc.wantOK {
				assert.False(
					t,
					ok,
					"Capture.Matches must report no match unless a plain run fills the pattern's first group",
				)
				return
			}
			require.True(t, ok)
			assert.Equal(
				t,
				tc.want,
				got,
				"Capture.Matches must return what the pattern's first group lifts out of the run",
			)
		})
	}
}

func TestItalic_Matches(t *testing.T) {
	cases := []struct {
		name   string
		inline prose.Inline
		want   string
		wantOK bool
	}{
		{
			name:   "returns the content of an emphasized run",
			inline: italic("article"),
			want:   "article",
			wantOK: true,
		},
		{
			name:   "returns the content verbatim, keeping the space around it",
			inline: italic(" True "),
			want:   " True ",
			wantOK: true,
		},
		{
			name:   "returns content written outside the ASCII range",
			inline: italic("Пра́вда"),
			want:   "Пра́вда",
			wantOK: true,
		},
		{
			name:   "returns the empty content of an emphasized run that carries none",
			inline: italic(""),
			want:   "",
			wantOK: true,
		},
		{
			name:   "returns nothing when the run carries no emphasis",
			inline: plain("article"),
			want:   "",
			wantOK: false,
		},
		{
			name:   "returns nothing when the run is set in bold",
			inline: prose.NewText("article", prose.StyleBold),
			want:   "",
			wantOK: false,
		},
		{
			name:   "returns nothing when the run is set in monowidth",
			inline: prose.NewText("article", prose.StyleCode),
			want:   "",
			wantOK: false,
		},
		{
			name:   "returns nothing when an emphasized link carries the content",
			inline: prose.NewLink("True", prose.StyleItalic, "#true"),
			want:   "",
			wantOK: false,
		},
		{
			name:   "returns nothing when the inline is a line break",
			inline: prose.NewLineBreak(),
			want:   "",
			wantOK: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := grammar.NewItalic().Matches(tc.inline)
			if !tc.wantOK {
				assert.False(
					t,
					ok,
					"Italic.Matches must report no match for anything but an emphasized text run",
				)
				return
			}
			require.True(t, ok)
			assert.Equal(
				t,
				tc.want,
				got,
				"Italic.Matches must return the content of the run verbatim, trimming nothing",
			)
		})
	}
}

func TestAnchor_Matches(t *testing.T) {
	cases := []struct {
		name   string
		inline prose.Inline
		want   string
		wantOK bool
	}{
		{
			name:   "returns the fragment an in-page link addresses, without its hash",
			inline: anchor("Message"),
			want:   "message",
			wantOK: true,
		},
		{
			name:   "returns the fragment of a link whose text does not name it",
			inline: prose.NewLink("these objects", prose.StylePlain, "#inlinequeryresult"),
			want:   "inlinequeryresult",
			wantOK: true,
		},
		{
			name:   "returns nothing when the link addresses an absolute URL",
			inline: prose.NewLink("guide", prose.StylePlain, "https://core.telegram.org/bots"),
			want:   "",
			wantOK: false,
		},
		{
			name:   "returns nothing when the link addresses a path from the site root",
			inline: prose.NewLink("payments", prose.StylePlain, "/bots/payments"),
			want:   "",
			wantOK: false,
		},
		{
			name:   "returns nothing when the inline is a plain text run",
			inline: plain("#message"),
			want:   "",
			wantOK: false,
		},
		{
			name:   "returns nothing when the inline is a line break",
			inline: prose.NewLineBreak(),
			want:   "",
			wantOK: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := grammar.NewAnchor().Matches(tc.inline)
			if !tc.wantOK {
				assert.False(
					t,
					ok,
					"Anchor.Matches must report no match for anything but a link into the same page",
				)
				return
			}
			require.True(t, ok)
			assert.Equal(
				t,
				tc.want,
				got,
				"Anchor.Matches must return the fragment the link addresses without its leading hash",
			)
		})
	}
}
