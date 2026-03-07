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

func newHTMLObjectField(t *testing.T, html string) parsing.HTMLObjectField {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	require.NoErrorf(t, err, "HTMLObjectField fixture must parse without error")
	return parsing.NewHTMLObjectField(dom.NewHTMLSelection(doc.Selection).Find("tr"))
}

func TestHTMLObjectField_Key(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "extracts snake_case key from first column",
			html: `<table><tr><td>message_id</td><td>Integer</td><td>Unique identifier</td></tr></table>`,
			want: "message_id",
		},
		{
			name:    "returns error when row has wrong number of columns",
			html:    `<table><tr><td>message_id</td><td>Integer</td></tr></table>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := newHTMLObjectField(t, tt.html).Key()
			if tt.wantErr {
				assert.Errorf(t, err, "HTMLObjectField must reject row with wrong column count")
				return
			}
			require.NoErrorf(t, err, "HTMLObjectField must extract key without error")
			got, err := key.Value()
			require.NoErrorf(t, err, "HTMLObjectField must produce a valid FieldKey")
			assert.Equalf(
				t,
				tt.want,
				got,
				"HTMLObjectField must extract key %q from first column",
				tt.want,
			)
		})
	}
}

func TestHTMLObjectField_Type(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		verify  func(t *testing.T, tree parsing.TypeTree)
		wantErr bool
	}{
		{
			name:   "extracts named type from second column",
			html:   `<table><tr><td>message_id</td><td>Integer</td><td>Unique identifier</td></tr></table>`,
			verify: assertHTMLObjectField_NamedType("Integer"),
		},
		{
			name:   "extracts array type from second column",
			html:   `<table><tr><td>photos</td><td>Array of PhotoSize</td><td>Available sizes</td></tr></table>`,
			verify: assertHTMLObjectField_ArrayType("PhotoSize"),
		},
		{
			name:    "returns error when row has wrong number of columns",
			html:    `<table><tr><td>message_id</td><td>Integer</td></tr></table>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := newHTMLObjectField(t, tt.html).Type()
			if tt.wantErr {
				assert.Errorf(t, err, "HTMLObjectField must reject row with wrong column count")
				return
			}
			require.NoErrorf(t, err, "HTMLObjectField must extract type without error")
			tt.verify(t, tree)
		})
	}
}

func TestHTMLObjectField_IsOptional(t *testing.T) {
	tests := []struct {
		name string
		html string
		want bool
	}{
		{
			name: "returns true when description starts with Optional",
			html: `<table><tr><td>username</td><td>String</td><td>Optional. User's username</td></tr></table>`,
			want: true,
		},
		{
			name: "returns false when description does not start with Optional",
			html: `<table><tr><td>message_id</td><td>Integer</td><td>Unique identifier</td></tr></table>`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHTMLObjectField(t, tt.html).IsOptional()
			require.NoErrorf(t, err, "HTMLObjectField must check optionality without error")
			assert.Equalf(
				t,
				tt.want,
				got,
				"HTMLObjectField must correctly detect optionality for %q",
				tt.name,
			)
		})
	}
}

func TestHTMLObjectField_Description(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "extracts description from third column",
			html: `<table><tr><td>message_id</td><td>Integer</td><td>Unique message identifier</td></tr></table>`,
			want: "Unique message identifier",
		},
		{
			name:    "returns error when row has wrong number of columns",
			html:    `<table><tr><td>message_id</td><td>Integer</td></tr></table>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHTMLObjectField(t, tt.html).Description()
			if tt.wantErr {
				assert.Errorf(t, err, "HTMLObjectField must reject row with wrong column count")
				return
			}
			require.NoErrorf(t, err, "HTMLObjectField must extract description without error")
			assert.Equalf(
				t,
				tt.want,
				got,
				"HTMLObjectField must extract description from third column",
			)
		})
	}
}

func assertHTMLObjectField_NamedType(want string) func(t *testing.T, tree parsing.TypeTree) {
	return func(t *testing.T, tree parsing.TypeTree) {
		t.Helper()
		root, err := tree.Root()
		require.NoErrorf(t, err, "HTMLObjectField must produce a parseable type tree")
		name, ok := root.Named()
		require.Truef(t, ok, "HTMLObjectField must produce a named type node")
		assert.Equalf(t, want, name, "HTMLObjectField must preserve type name %q", want)
	}
}

func assertHTMLObjectField_ArrayType(wantInner string) func(t *testing.T, tree parsing.TypeTree) {
	return func(t *testing.T, tree parsing.TypeTree) {
		t.Helper()
		root, err := tree.Root()
		require.NoErrorf(t, err, "HTMLObjectField must produce a parseable type tree")
		inner, ok := root.Array()
		require.Truef(t, ok, "HTMLObjectField must produce an array type node")
		name, ok := inner.Named()
		require.Truef(t, ok, "HTMLObjectField must produce a named node as array element")
		assert.Equalf(
			t,
			wantInner,
			name,
			"HTMLObjectField must preserve array element type %q",
			wantInner,
		)
	}
}
