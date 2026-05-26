// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model/spec/gq"
	pkggq "github.com/andreychh/tgen/pkg/gq"
)

func tdWith(inner string) pkggq.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(
		`<table><tbody><tr><td>` + inner + `</td></tr></tbody></table>`,
	))
	return pkggq.NewNormSelection(doc.Find("td"))
}

func emptySel() pkggq.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<div></div>`))
	return pkggq.NewNormSelection(doc.Find("td"))
}

func TestObjectFieldDescription_Value(t *testing.T) {
	cases := []struct {
		name    string
		sel     pkggq.Selection
		want    string
		wantErr bool
	}{
		{
			name: "strips Optional. prefix from an optional field description",
			sel:  tdWith("Optional. The message was edited."),
			want: "The message was edited.",
		},
		{
			name: "returns description unchanged when the Optional. prefix is absent",
			sel:  tdWith("Unique identifier of the target chat."),
			want: "Unique identifier of the target chat.",
		},
		{
			name: "returns description unchanged when Optional. appears mid-text rather than as prefix",
			sel:  tdWith("This field is Optional. when not provided."),
			want: "This field is Optional. when not provided.",
		},
		{
			name:    "returns error when the selection is empty",
			sel:     emptySel(),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gq.NewObjectFieldDescription(tc.sel).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ObjectFieldDescription.Value must return an error when the td selection is empty",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ObjectFieldDescription.Value must strip only the leading 'Optional. ' prefix and return the remainder as the description",
			)
		})
	}
}

func TestObjectFieldDescription_Links(t *testing.T) {
	cases := []struct {
		name    string
		sel     pkggq.Selection
		want    []string
		wantErr bool
	}{
		{
			name: "returns all hrefs from multiple anchor tags in the description",
			sel: tdWith(
				`See <a href="#sendmessage">sendMessage</a> or <a href="#getchat">getChat</a>`,
			),
			want: []string{"#sendmessage", "#getchat"},
		},
		{
			name: "returns nil when the description contains no anchor tags",
			sel:  tdWith("Plain text description without links."),
			want: nil,
		},
		{
			name:    "returns error when the selection is empty",
			sel:     emptySel(),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gq.NewObjectFieldDescription(tc.sel).Links()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ObjectFieldDescription.Links must return an error when the td selection is empty",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ObjectFieldDescription.Links must return the hrefs of all anchor tags found in the description",
			)
		})
	}
}
