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

func TestRawObject_ID(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		wantID  string
		wantErr bool
	}{
		{
			name:   "returns valid ID from anchor",
			html:   `<h4><a class="anchor" href="#message"></a>Message</h4>`,
			wantID: "#message",
		},
		{
			name:    "returns error when anchor is missing",
			html:    `<h4>Message</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error when href is missing",
			html:    `<h4><a class="anchor"></a>Message</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error when format is invalid (no hash)",
			html:    `<h4><a class="anchor" href="message"></a>Message</h4>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawObject(
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

func TestRawObject_Name(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantName string
		wantErr  bool
	}{
		{
			name:     "returns valid PascalCase name",
			html:     `<h4><a class="anchor" href="#message"></a>Message</h4>`,
			wantName: "Message",
		},
		{
			name:    "returns error for invalid casing (camelCase)",
			html:    `<h4>message</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error for invalid characters",
			html:    `<h4>Message_1</h4>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawObject(
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

func TestRawObject_Description(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantDesc string
	}{
		{
			name: "extracts description paragraphs before table",
			html: `
				<h4>Message</h4>
				<p>This is a message object.</p>
				<p>It contains data.</p>
				<table><tbody><tr><td>...</td></tr></tbody></table>
			`,
			wantDesc: "This is a message object. It contains data.",
		},
		{
			name: "stops extraction at next header if table is missing",
			html: `
				<h4>Message</h4>
				<p>Description only.</p>
				<h4>NextObject</h4>
			`,
			wantDesc: "Description only.",
		},
		{
			name: "returns empty string if no paragraphs exist",
			html: `
				<h4>Message</h4>
				<table><tbody><tr><td>...</td></tr></tbody></table>
			`,
			wantDesc: "",
		},
		{
			name: "does not include content after the table",
			html: `
				<h4>Message</h4>
				<p>Valid description.</p>
				<table><tbody><tr><td>...</td></tr></tbody></table>
				<p>Note: this paragraph belongs to the fields or is a generic note.</p>
			`,
			wantDesc: "Valid description.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawObject(
				dom.NewHTMLSelection(doc.Selection).Find("h4").First(),
			).Description()
			require.NoError(t, err, "unexpected error returned during description extraction")
			assert.Equal(t, tt.wantDesc, got, "extracted description does not match expectation")
		})
	}
}

func TestRawObject_Fields(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantKeys []string
	}{
		{
			name: "extracts all fields from the object's table",
			html: `
				<h4>Message</h4>
				<p>Desc</p>
				<table>
					<tbody>
						<tr><td>message_id</td><td>Integer</td><td>UniqueId</td></tr>
						<tr><td>date</td><td>Integer</td><td>DateSent</td></tr>
					</tbody>
				</table>
				<h4>NextObject</h4>
			`,
			wantKeys: []string{"message_id", "date"},
		},
		{
			name: "returns empty sequence if table is missing",
			html: `
				<h4>Message</h4>
				<p>No fields for this object.</p>
				<h4>NextObject</h4>
			`,
			wantKeys: []string{},
		},
		{
			name: "ignores tables belonging to the next section",
			html: `
				<h4>Message</h4>
				<p>Desc</p>
				<h4>NextObject</h4>
				<table>
					<tbody>
						<tr><td>update_id</td><td>Integer</td><td>UpdateId</td></tr>
					</tbody>
				</table>
			`,
			wantKeys: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			obj := parsing.NewRawObject(
				dom.NewHTMLSelection(doc.Selection).Find("h4").First(),
			)
			keys := []string{}
			for field := range obj.Fields() {
				key, err := field.JSONKey()
				require.NoError(t, err, "field JSONKey extraction failed inside iterator")
				keys = append(keys, key)
			}
			assert.Equal(
				t,
				len(tt.wantKeys),
				len(keys),
				"number of extracted fields does not match",
			)
			assert.Equal(t, tt.wantKeys, keys, "extracted JSON keys do not match expectation")
		})
	}
}
