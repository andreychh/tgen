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

func TestRawVariant_ID(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		wantID  string
		wantErr bool
	}{
		{
			name:   "returns valid ID from href",
			html:   `<li><a href="#message">Message</a></li>`,
			wantID: "#message",
		},
		{
			name:    "returns error when anchor is missing",
			html:    `<li>Message</li>`,
			wantErr: true,
		},
		{
			name:    "returns error when href is missing",
			html:    `<li><a>Message</a></li>`,
			wantErr: true,
		},
		{
			name:    "returns error when format is invalid (no hash)",
			html:    `<li><a href="message">Message</a></li>`,
			wantErr: true,
		},
		{
			name:    "returns error for empty href",
			html:    `<li><a href="">Message</a></li>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewDefaultRawVariant(
				dom.NewHTMLSelection(doc.Selection).Find("li"),
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

func TestRawVariant_Name(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantName string
		wantErr  bool
	}{
		{
			name:     "returns valid name",
			html:     `<li><a href="#message">Message</a></li>`,
			wantName: "Message",
		},
		{
			name:     "returns valid name with whitespace (trimmed)",
			html:     `<li>   User   </li>`,
			wantName: "User",
		},
		{
			name:    "returns error for invalid casing (camelCase)",
			html:    `<li>message</li>`,
			wantErr: true,
		},
		{
			name:    "returns error for invalid characters (underscore)",
			html:    `<li>Message_Type</li>`,
			wantErr: true,
		},
		{
			name:    "returns error for empty text",
			html:    `<li></li>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewDefaultRawVariant(
				dom.NewHTMLSelection(doc.Selection).Find("li"),
			).Name()
			if tt.wantErr {
				assert.Error(
					t,
					err,
					"validation did not return expected error for invalid name format",
				)
				return
			}
			require.NoError(t, err, "unexpected error returned during name extraction")
			assert.Equal(t, tt.wantName, got, "extracted name does not match expectation")
		})
	}
}
