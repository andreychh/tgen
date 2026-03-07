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

func newHTMLObject(t *testing.T, html string) parsing.HTMLObject {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	require.NoErrorf(t, err, "HTMLObject fixture must parse without error")
	return parsing.NewHTMLObject(dom.NewHTMLSelection(doc.Selection).Find("h4"))
}

func TestHTMLObject_Ref(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "extracts ref without hash prefix from anchor href",
			html: `<div><h4><a class="anchor" href="#message">Message</a></h4></div>`,
			want: "message",
		},
		{
			name:    "returns error when anchor element is missing",
			html:    `<div><h4>Message</h4></div>`,
			wantErr: true,
		},
		{
			name:    "returns error when anchor href attribute is missing",
			html:    `<div><h4><a class="anchor">Message</a></h4></div>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := newHTMLObject(t, tt.html).Ref()
			if tt.wantErr {
				assert.Errorf(t, err, "HTMLObject must reject header without valid anchor href")
				return
			}
			require.NoErrorf(t, err, "HTMLObject must extract ref without error")
			got, err := ref.Value()
			require.NoErrorf(t, err, "HTMLObject must produce a valid DefinitionRef")
			assert.Equalf(t, tt.want, got, "HTMLObject must strip hash and return ref %q", tt.want)
		})
	}
}

func TestHTMLObject_Name(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "extracts PascalCase name from header text",
			html: `<div><h4><a class="anchor" href="#send-message">SendMessage</a></h4></div>`,
			want: "SendMessage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := newHTMLObject(t, tt.html)
			n, err := obj.Name()
			require.NoErrorf(t, err, "HTMLObject must extract name without error")
			got, err := n.Value()
			require.NoErrorf(t, err, "HTMLObject must produce a valid ObjectName")
			assert.Equalf(t, tt.want, got, "HTMLObject must extract name %q from header text", tt.want)
		})
	}
}

func TestHTMLObject_Description(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "extracts single paragraph description",
			html: `<div><h4><a class="anchor" href="#message">Message</a></h4><p>This object represents a message.</p></div>`,
			want: "This object represents a message.",
		},
		{
			name: "joins multiple paragraphs with space",
			html: `<div><h4><a class="anchor" href="#message">Message</a></h4><p>First paragraph.</p><p>Second paragraph.</p></div>`,
			want: "First paragraph. Second paragraph.",
		},
		{
			name: "returns empty string when no paragraphs follow header",
			html: `<div><h4><a class="anchor" href="#message">Message</a></h4></div>`,
			want: "",
		},
		{
			name: "stops at next header",
			html: `<div><h4><a class="anchor" href="#message">Message</a></h4><p>Description.</p><h4><a class="anchor" href="#other">Other</a></h4><p>Other description.</p></div>`,
			want: "Description.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHTMLObject(t, tt.html).Description()
			require.NoErrorf(t, err, "HTMLObject must extract description without error")
			assert.Equalf(t, tt.want, got, "HTMLObject must produce description %q", tt.want)
		})
	}
}

func TestHTMLObject_Fields(t *testing.T) {
	tests := []struct {
		name   string
		html   string
		verify func(t *testing.T, fields []parsing.Field)
	}{
		{
			name: "yields all rows from object field table",
			html: `<div>
				<h4><a class="anchor" href="#message">Message</a></h4>
				<p>A message.</p>
				<table><tbody>
					<tr><td>message_id</td><td>Integer</td><td>Unique identifier</td></tr>
					<tr><td>text</td><td>String</td><td>Message text</td></tr>
				</tbody></table>
			</div>`,
			verify: assertHTMLObject_FieldCount(2),
		},
		{
			name:   "yields no fields when table is absent",
			html:   `<div><h4><a class="anchor" href="#message">Message</a></h4><p>A message.</p></div>`,
			verify: assertHTMLObject_FieldCount(0),
		},
		{
			name: "stops at next header",
			html: `<div>
				<h4><a class="anchor" href="#message">Message</a></h4>
				<table><tbody>
					<tr><td>message_id</td><td>Integer</td><td>Unique identifier</td></tr>
				</tbody></table>
				<h4><a class="anchor" href="#other">Other</a></h4>
				<table><tbody>
					<tr><td>other_id</td><td>Integer</td><td>Other identifier</td></tr>
					<tr><td>name</td><td>String</td><td>Name</td></tr>
				</tbody></table>
			</div>`,
			verify: assertHTMLObject_FieldCount(1),
		},
		{
			name: "yields field with correct key from first column",
			html: `<div>
				<h4><a class="anchor" href="#message">Message</a></h4>
				<table><tbody>
					<tr><td>chat_id</td><td>Integer or String</td><td>Target chat</td></tr>
				</tbody></table>
			</div>`,
			verify: assertHTMLObject_FirstFieldKey("chat_id"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fields []parsing.Field
			for f := range newHTMLObject(t, tt.html).Fields() {
				fields = append(fields, f)
			}
			tt.verify(t, fields)
		})
	}
}

func assertHTMLObject_FieldCount(want int) func(t *testing.T, fields []parsing.Field) {
	return func(t *testing.T, fields []parsing.Field) {
		t.Helper()
		assert.Equalf(t, want, len(fields), "HTMLObject must yield exactly %d field(s)", want)
	}
}

func assertHTMLObject_FirstFieldKey(want string) func(t *testing.T, fields []parsing.Field) {
	return func(t *testing.T, fields []parsing.Field) {
		t.Helper()
		require.NotEmptyf(t, fields, "HTMLObject must yield at least one field")
		key, err := fields[0].Key()
		require.NoErrorf(t, err, "HTMLObject must produce a field with extractable key")
		got, err := key.Value()
		require.NoErrorf(t, err, "HTMLObject must produce a field with valid key")
		assert.Equalf(t, want, got, "HTMLObject must yield first field with key %q", want)
	}
}
