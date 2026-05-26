// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq_test

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec/gq"
	pkggq "github.com/andreychh/tgen/pkg/gq"
)

func h4With(text string) pkggq.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<h4>` + text + `</h4>`))
	return pkggq.NewNormSelection(doc.Find("h4"))
}

func aWith(href string) pkggq.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<a href="` + href + `">link</a>`))
	return pkggq.NewNormSelection(doc.Find("a"))
}

func strongWith(text string) pkggq.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<strong>` + text + `</strong>`))
	return pkggq.NewNormSelection(doc.Find("strong"))
}

func TestName_Value(t *testing.T) {
	cases := []struct {
		name    string
		sel     pkggq.Selection
		want    model.Name
		wantErr bool
	}{
		{
			name: "returns the definition name from an h4 with alphanumeric text",
			sel:  h4With("SendMessage"),
			want: "SendMessage",
		},
		{
			name: "returns the definition name when the text mixes letters and digits",
			sel:  h4With("getMe2"),
			want: "getMe2",
		},
		{
			name:    "returns error when the name contains an underscore",
			sel:     h4With("send_message"),
			wantErr: true,
		},
		{
			name:    "returns error when the name contains non-ASCII characters",
			sel:     h4With("Ñoño"),
			wantErr: true,
		},
		{
			name:    "returns error when the selection is empty",
			sel:     emptySel(),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gq.NewName(tc.sel).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Name.Value must return an error when the selection is empty or the name contains non-alphanumeric characters",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Name.Value must return the h4 text as the definition name",
			)
		})
	}
}

func TestKey_Value(t *testing.T) {
	cases := []struct {
		name    string
		sel     pkggq.Selection
		want    model.Key
		wantErr bool
	}{
		{
			name: "returns the field key from a td with a snake_case identifier",
			sel:  tdWith("message_id"),
			want: "message_id",
		},
		{
			name: "returns the field key from a td with a single lowercase word",
			sel:  tdWith("chat"),
			want: "chat",
		},
		{
			name:    "returns error when the key starts with an uppercase letter",
			sel:     tdWith("ChatId"),
			wantErr: true,
		},
		{
			name:    "returns error when the key starts with a digit",
			sel:     tdWith("2fa_enabled"),
			wantErr: true,
		},
		{
			name:    "returns error when the selection is empty",
			sel:     emptySel(),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gq.NewKey(tc.sel).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"Key.Value must return an error when the selection is empty or the key does not match the lowercase identifier pattern",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"Key.Value must return the td text as the field key",
			)
		})
	}
}

func TestReleaseDate_Value(t *testing.T) {
	parsed, _ := time.Parse("#January-2-2006", "#may-8-2026")
	cases := []struct {
		name    string
		sel     pkggq.Selection
		want    model.ReleaseDate
		wantErr bool
	}{
		{
			name: "parses release date from a valid anchor href",
			sel:  aWith("#may-8-2026"),
			want: model.ReleaseDate(parsed),
		},
		{
			name:    "returns error when the href does not match the date format",
			sel:     aWith("#sendmessage"),
			wantErr: true,
		},
		{
			name:    "returns error when the selection is empty",
			sel:     emptySel(),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gq.NewReleaseDate(tc.sel).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ReleaseDate.Value must return an error when the selection is empty or the href does not match the '#Month-day-year' format",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ReleaseDate.Value must parse the release date from the anchor href",
			)
		})
	}
}

func TestReleaseVersion_Value(t *testing.T) {
	cases := []struct {
		name    string
		sel     pkggq.Selection
		want    model.ReleaseVersion
		wantErr bool
	}{
		{
			name: "extracts version number from Bot API X.Y text",
			sel:  strongWith("Bot API 10.0"),
			want: "10.0",
		},
		{
			name: "extracts version number with minor version greater than zero",
			sel:  strongWith("Bot API 7.4"),
			want: "7.4",
		},
		{
			name:    "returns error when the version number has no decimal part",
			sel:     strongWith("Bot API 10"),
			wantErr: true,
		},
		{
			name:    "returns error when the text is missing the Bot API prefix",
			sel:     strongWith("API 10.0"),
			wantErr: true,
		},
		{
			name:    "returns error when the selection is empty",
			sel:     emptySel(),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gq.NewReleaseVersion(tc.sel).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ReleaseVersion.Value must return an error when the selection is empty or the text does not match 'Bot API X.Y'",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ReleaseVersion.Value must extract the X.Y version number from the 'Bot API X.Y' text",
			)
		})
	}
}