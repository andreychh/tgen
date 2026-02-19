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

func TestRawSpecification_Unions(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		limit     int
		wantNames []string
	}{
		{
			name: "extracts only union definitions from a document containing mixed entities.",
			html: `
				<h4><a class="anchor" href="#getupdates"></a>getUpdates</h4>
				<p>Method desc</p>
				<table><tr><td>offset</td></tr></table>

				<h4><a class="anchor" href="#user"></a>User</h4>
				<p>Object desc</p>
				<table><tr><td>id</td></tr></table>

				<h4><a class="anchor" href="#maybeinaccessiblemessage"></a>MaybeInaccessibleMessage</h4>
				<p>Union 1 desc</p>
				<ul><li>Message</li><li>InaccessibleMessage</li></ul>

				<h4><a class="anchor" href="#messageorigin"></a>MessageOrigin</h4>
				<p>Union 2 desc</p>
				<ul><li>MessageOriginUser</li><li>MessageOriginChat</li></ul>
			`,
			limit:     -1,
			wantNames: []string{"MaybeInaccessibleMessage", "MessageOrigin"},
		},
		{
			name: "yields an empty sequence when the document contains no union definitions.",
			html: `
				<h4><a class="anchor" href="#getupdates"></a>getUpdates</h4>
				<table><tr><td>offset</td></tr></table>
				<h4><a class="anchor" href="#user"></a>User</h4>
				<table><tr><td>id</td></tr></table>
			`,
			limit:     -1,
			wantNames: []string{},
		},
		{
			name: "stops iterating immediately when the yield function returns false.",
			html: `
				<h4><a class="anchor" href="#unionone"></a>UnionOne</h4>
				<ul><li>VariantA</li></ul>
				<h4><a class="anchor" href="#uniontwo"></a>UnionTwo</h4>
				<ul><li>VariantB</li></ul>
				<h4><a class="anchor" href="#unionthree"></a>UnionThree</h4>
				<ul><li>VariantC</li></ul>
			`,
			limit:     1,
			wantNames: []string{"UnionOne"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			spec := parsing.NewRawSpecification(dom.NewHTMLSelection(doc.Selection))
			gotNames := []string{}
			count := 0
			for union := range spec.Unions() {
				name, err := union.Name()
				require.NoError(t, err, "extracted union name is invalid or missing")
				gotNames = append(gotNames, name)
				count++
				if tt.limit != -1 && count >= tt.limit {
					break
				}
			}
			assert.Equal(
				t,
				tt.wantNames,
				gotNames,
				"extracted union sequence does not match expected items",
			)
		})
	}
}
