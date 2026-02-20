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

func TestRawTypeField_JSONKey(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "extracts valid snake_case key",
			html: `<table><tr><td>message_id</td><td>Integer</td><td>UniqueId</td></tr></table>`,
			want: "message_id",
		},
		{
			name:    "returns error for invalid characters",
			html:    `<table><tr><td>Message-Id</td><td>Integer</td><td>Desc</td></tr></table>`,
			wantErr: true,
		},
		{
			name:    "returns error when columns are missing",
			html:    `<table><tr><td>message_id</td><td>Integer</td></tr></table>`,
			wantErr: true,
		},
		{
			name: "extracts valid key containing numbers",
			html: `<table><tr><td>mpeg4_url</td><td>String</td><td>A valid URL for the MPEG4 file</td></tr></table>`,
			want: "mpeg4_url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawObjectField(
				dom.NewHTMLSelection(doc.Selection).Find("tr"),
			).JSONKey()
			if tt.wantErr {
				assert.Error(
					t,
					err,
					"validation did not return expected error for invalid JSON key",
				)
				return
			}
			require.NoError(t, err, "unexpected error returned during JSON key extraction")
			assert.Equal(t, tt.want, got, "extracted JSON key does not match expectation")
		})
	}
}

func TestRawTypeField_Type(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "extracts simple type",
			html: `<table><tr><td>id</td><td>Integer</td><td>Desc</td></tr></table>`,
			want: "Integer",
		},
		{
			name: "extracts complex array type",
			html: `<table><tr><td>users</td><td>Array of User</td><td>Desc</td></tr></table>`,
			want: "Array of User",
		},
		{
			name: "extracts union type",
			html: `<table><tr><td>chat_id</td><td>Integer or String</td><td>Desc</td></tr></table>`,
			want: "Integer or String",
		},
		{
			name:    "returns error for invalid type format",
			html:    `<table><tr><td>id</td><td>Integer!</td><td>Desc</td></tr></table>`,
			wantErr: true,
		},
		{
			name:    "returns error when columns are missing",
			html:    `<table><tr><td>message_id</td><td>Integer</td></tr></table>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawObjectField(
				dom.NewHTMLSelection(doc.Selection).Find("tr"),
			).Type()
			if tt.wantErr {
				assert.Error(
					t,
					err,
					"validation did not return expected error for invalid type format",
				)
				return
			}
			require.NoError(t, err, "unexpected error returned during type extraction")
			assert.Equal(t, tt.want, got, "extracted type does not match expectation")
		})
	}
}

func TestRawTypeField_Description(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "extracts full description text",
			html: `<table><tr><td>key</td><td>Type</td><td>This is a long description.</td></tr></table>`,
			want: "This is a long description.",
		},
		{
			name:    "returns error when columns are missing",
			html:    `<table><tr><td>key</td><td>Type</td></tr></table>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawObjectField(
				dom.NewHTMLSelection(doc.Selection).Find("tr"),
			).Description()
			if tt.wantErr {
				assert.Error(t, err, "validation did not return expected error for missing columns")
				return
			}
			require.NoError(t, err, "unexpected error returned during description extraction")
			assert.Equal(t, tt.want, got, "extracted description does not match expectation")
		})
	}
}

func TestRawTypeField_IsOptional(t *testing.T) {
	tests := []struct {
		name string
		html string
		want bool
	}{
		{
			name: "returns true for Optional prefix",
			html: `<table><tr><td>key</td><td>Type</td><td>Optional. Description.
</td></tr></table>`,
			want: true,
		},
		{
			name: "returns false for mandatory fields",
			html: `<table><tr><td>key</td><td>Type</td><td>Description.</td></tr></table>`,
			want: false,
		},
		{
			name: "returns false if Optional is in the middle",
			html: `<table><tr><td>key</td><td>Type</td><td>This field is Optional.
</td></tr></table>`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawObjectField(
				dom.NewHTMLSelection(doc.Selection).Find("tr"),
			).IsOptional()
			require.NoError(t, err, "unexpected error returned during optionality check")
			assert.Equal(t, tt.want, got, "optionality status does not match expectation")
		})
	}
}

func TestRawTypeField_Name(t *testing.T) {
	html := `<table><tr><td>message_id</td><td>Integer</td><td>Desc</td></tr></table>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	require.NoError(t, err, "HTML fixture does not parse correctly")
	got, err := parsing.NewRawObjectField(
		dom.NewHTMLSelection(doc.Selection).Find("tr"),
	).Name()
	require.NoError(t, err, "unexpected error returned during name extraction")
	assert.Equal(t, "message_id", got, "extracted name does not match expected JSON key delegation")
}
