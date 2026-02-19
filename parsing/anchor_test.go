// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/dom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnchor_Kind(t *testing.T) {
	tests := []struct {
		name string
		html string
		want parsing.DefinitionKind
	}{
		{
			name: "identifies a method by its lowercase header",
			html: `
				<h4><a class="anchor" href="#getupdates"></a>getUpdates</h4>
				<p>Use this method to receive incoming updates.</p>
				<table><tr><td>offset</td></tr></table>
			`,
			want: parsing.KindMethod,
		},
		{
			name: "identifies an object by an uppercase header and a table",
			html: `
				<h4><a class="anchor" href="#user"></a>User</h4>
				<p>This object represents a Telegram user or bot.</p>
				<table><tr><td>id</td></tr></table>
			`,
			want: parsing.KindObject,
		},
		{
			name: "identifies a union by an uppercase header and a list without a table",
			html: `
				<h4><a class="anchor" href="#maybeinaccessiblemessage"></a>MaybeInaccessibleMessage</h4>
				<p>This object describes a message that can be inaccessible.</p>
				<ul>
					<li>Message</li>
					<li>InaccessibleMessage</li>
				</ul>
			`,
			want: parsing.KindUnion,
		},
		{
			name: "returns unknown when the anchor href is missing",
			html: `
				<h4><a class="anchor"></a>User</h4>
				<table><tr><td>id</td></tr></table>
			`,
			want: parsing.KindUnknown,
		},
		{
			name: "returns unknown when the anchor href contains a hyphen",
			html: `
				<h4><a class="anchor" href="#formatting-options"></a>Formatting Options</h4>
				<p>The Bot API supports basic formatting...</p>
			`,
			want: parsing.KindUnknown,
		},
		{
			name: "returns unknown when the section lacks both table and list",
			html: `
				<h4><a class="anchor" href="#gettingupdates"></a>Getting updates</h4>
				<p>There are two mutually exclusive ways of receiving updates...</p>
			`,
			want: parsing.KindUnknown,
		},
		{
			name: "stops scanning at the next section header",
			html: `
				<h4><a class="anchor" href="#user"></a>User</h4>
				<p>User description.</p>
				<table><tr><td>id</td></tr></table>
				<h4><a class="anchor" href="#chat"></a>Chat</h4>
				<ul><li>Wrong list belonging to Chat</li></ul>
			`,
			want: parsing.KindObject,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			anchor := parsing.NewAnchor(dom.NewHTMLSelection(doc.Find("h4").First()))
			assert.Equal(
				t,
				tt.want,
				anchor.Kind(),
				"extracted definition kind does not match expectation",
			)
		})
	}
}
