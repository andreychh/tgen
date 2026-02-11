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

func TestRawUnion_ID(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		wantID  string
		wantErr bool
	}{
		{
			name:   "returns valid ID from anchor",
			html:   `<h4><a class="anchor" href="#uniontype"></a>UnionType</h4>`,
			wantID: "#uniontype",
		},
		{
			name:   "returns valid ID with whitespace (trimmed)",
			html:   `<h4><a class="anchor" href="  #uniontype  "></a>UnionType</h4>`,
			wantID: "#uniontype",
		},
		{
			name:    "returns error when anchor is missing",
			html:    `<h4>UnionType</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error when href is missing",
			html:    `<h4><a class="anchor"></a>UnionType</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error when format is invalid (no hash)",
			html:    `<h4><a class="anchor" href="uniontype"></a>UnionType</h4>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewDefaultRawUnion(
				dom.NewHTMLSelection(doc.Selection).Find("h4"),
			).ID()
			if tt.wantErr {
				assert.Error(t, err, "validation did not return expected error for invalid input")
				return
			}
			require.NoError(t, err, "unexpected error returned during extraction")
			assert.Equal(t, tt.wantID, got, "extracted ID does not match expectation")
		})
	}
}

func TestRawUnion_Name(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantName string
		wantErr  bool
	}{
		{
			name:     "returns valid PascalCase name",
			html:     `<h4><a class="anchor" href="#union"></a>UnionType</h4>`,
			wantName: "UnionType",
		},
		{
			name:     "returns normalized name (whitespace removal)",
			html:     `<h4>   UnionType   </h4>`,
			wantName: "UnionType",
		},
		{
			name:    "returns error for invalid casing (camelCase)",
			html:    `<h4>unionType</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error for invalid characters",
			html:    `<h4>Union_Type</h4>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewDefaultRawUnion(
				dom.NewHTMLSelection(doc.Selection).Find("h4"),
			).Name()
			if tt.wantErr {
				assert.Error(t, err, "validation did not return expected error for invalid name")
				return
			}
			require.NoError(t, err, "unexpected error returned during name extraction")
			assert.Equal(t, tt.wantName, got, "extracted name does not match expectation")
		})
	}
}

func TestRawUnion_Description(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantDesc string
	}{
		{
			name: "extracts description paragraphs before list",
			html: `
				<h4>Union</h4>
				<p>First paragraph.</p>
				<p>Second paragraph.</p>
				<ul><li>Variant</li></ul>
			`,
			wantDesc: "First paragraph. Second paragraph.",
		},
		{
			name: "stops extraction at next header if list is missing",
			html: `
				<h4>Union</h4>
				<p>Description.</p>
				<h4>NextUnion</h4>
			`,
			wantDesc: "Description.",
		},
		{
			name: "returns empty string if no paragraphs exist",
			html: `
				<h4>Union</h4>
				<ul><li>Variant</li></ul>
			`,
			wantDesc: "",
		},
		{
			name: "does not include content after the list",
			html: `
				<h4>Union</h4>
				<p>Description.</p>
				<ul><li>Variant</li></ul>
				<p>Note: this is not part of description.</p>
			`,
			wantDesc: "Description.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewDefaultRawUnion(
				dom.NewHTMLSelection(doc.Selection).Find("h4").First(),
			).Description()
			require.NoError(t, err, "unexpected error returned during description extraction")
			assert.Equal(t, tt.wantDesc, got, "extracted description does not match expectation")
		})
	}
}

func TestRawUnion_Variants(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		wantNames []string
	}{
		{
			name: "extracts all variants from the list",
			html: `
				<h4>Union</h4>
				<p>Desc</p>
				<ul>
					<li><a href="#a">VariantA</a></li>
					<li><a href="#b">VariantB</a></li>
				</ul>
				<h4>Next</h4>
			`,
			wantNames: []string{"VariantA", "VariantB"},
		},
		{
			name: "returns empty sequence if list is missing",
			html: `
				<h4>Union</h4>
				<p>Desc</p>
				<h4>Next</h4>
			`,
			wantNames: []string{},
		},
		{
			name: "ignores lists belonging to the next section",
			html: `
				<h4>Union</h4>
				<p>Desc</p>
				<h4>NextUnion</h4>
				<ul>
					<li><a href="#c">VariantC</a></li>
				</ul>
			`,
			wantNames: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			union := parsing.NewDefaultRawUnion(
				dom.NewHTMLSelection(doc.Selection).Find("h4").First(),
			)
			names := []string{}
			for v := range union.Variants() {
				name, err := v.Name()
				require.NoError(t, err, "variant name extraction failed inside iterator")
				names = append(names, name)
			}
			assert.Equal(
				t,
				len(tt.wantNames),
				len(names),
				"number of extracted variants does not match",
			)
			assert.Equal(t, tt.wantNames, names, "variant names do not match expectation")
		})
	}
}
