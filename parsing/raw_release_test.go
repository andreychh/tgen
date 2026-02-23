// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing_test

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/dom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawRelease_ID(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		wantID  string
		wantErr bool
	}{
		{
			name:   "returns valid ID from anchor",
			html:   `<h4><a class="anchor" href="#february-9-2026"></a>February 9, 2026</h4>`,
			wantID: "#february-9-2026",
		},
		{
			name:    "returns error when anchor is missing",
			html:    `<h4>February 9, 2026</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error when href is missing",
			html:    `<h4><a class="anchor"></a>February 9, 2026</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error when format is invalid",
			html:    `<h4><a class="anchor" href="#invalid_Release!"></a>February 9, 2026</h4>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawRelease(
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

func TestRawRelease_Version(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		wantVersion string
		wantErr     bool
	}{
		{
			name: "extracts valid version from next paragraph",
			html: `
				<h4><a class="anchor" href="#february-9-2026"></a>February 9, 2026</h4>
				<p><strong>Bot API 9.4</strong></p>
			`,
			wantVersion: "v9.4",
		},
		{
			name: "returns error if strong tag is missing",
			html: `
				<h4>February 9, 2026</h4>
				<p>Bot API 9.4</p>
			`,
			wantErr: true,
		},
		{
			name: "returns error if text does not match version pattern",
			html: `
				<h4>February 9, 2026</h4>
				<p><strong>Some other update</strong></p>
			`,
			wantErr: true,
		},
		{
			name: "returns error if no paragraphs follow the header",
			html: `
             <h4><a class="anchor" href="#february-9-2026"></a>February 9, 2026</h4>
             <div><strong>Bot API 9.4</strong> inside wrong tag</div>
          `,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawRelease(
				dom.NewHTMLSelection(doc.Selection).Find("h4"),
			).Version()
			if tt.wantErr {
				assert.Error(
					t,
					err,
					"validation did not return expected error for invalid version text",
				)
				return
			}
			require.NoError(t, err, "unexpected error returned during version extraction")
			assert.Equal(t, tt.wantVersion, got, "extracted version does not match expectation")
		})
	}
}

func TestRawRelease_Date(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantDate time.Time
		wantErr  bool
	}{
		{
			name:     "extracts valid date from anchor href",
			html:     `<h4><a class="anchor" href="#february-9-2026"></a>February 9, 2026</h4>`,
			wantDate: time.Date(2026, time.February, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "returns error when anchor is missing",
			html:    `<h4>February 9, 2026</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error when href attribute is missing",
			html:    `<h4><a class="anchor"></a>February 9, 2026</h4>`,
			wantErr: true,
		},
		{
			name:    "returns error for invalid date format in href",
			html:    `<h4><a class="anchor" href="#not-a-date"></a>February 9, 2026</h4>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got, err := parsing.NewRawRelease(
				dom.NewHTMLSelection(doc.Selection).Find("h4"),
			).Date()
			if tt.wantErr {
				assert.Error(t, err, "validation did not return expected error for invalid input")
				return
			}
			require.NoError(t, err, "unexpected error returned during date extraction")
			assert.Equal(t, tt.wantDate, got, "extracted date does not match expectation")
		})
	}
}
