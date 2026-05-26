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
	"github.com/andreychh/tgen/model/types"
	pkggq "github.com/andreychh/tgen/pkg/gq"
)

// returnTypeFixture builds a document with a method h4 (getMe) whose body
// contains a single <p> with the given text, followed by Message and Chat
// object h4s so that Catalog.Lookup can resolve those names.
func returnTypeFixture(desc string) (pkggq.Selection, pkggq.Selection) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(
		`<div id="dev_page_content">` +
			`<h4><a class="anchor" href="#getme">getMe</a></h4>` +
			`<p>` + desc + `</p>` +
			`<h4><a class="anchor" href="#message">Message</a></h4>` +
			`<table><tbody><tr><td>id</td><td>Integer</td><td>The ID.</td></tr></tbody></table>` +
			`<h4><a class="anchor" href="#chat">Chat</a></h4>` +
			`<table><tbody><tr><td>id</td><td>Integer</td><td>The ID.</td></tr></tbody></table>` +
			`</div>`,
	))
	return pkggq.NewNormSelection(doc.Selection), pkggq.NewNormSelection(doc.Find("h4").First())
}

// returnTypeNoParaFixture builds a document whose method h4 has no <p>
// siblings before the next h4, so ReturnType.Value returns an error.
func returnTypeNoParaFixture() (pkggq.Selection, pkggq.Selection) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(
		`<div id="dev_page_content">` +
			`<h4><a class="anchor" href="#getme">getMe</a></h4>` +
			`<h4><a class="anchor" href="#message">Message</a></h4>` +
			`<table><tbody><tr><td>id</td><td>Integer</td><td>The ID.</td></tr></tbody></table>` +
			`</div>`,
	))
	return pkggq.NewNormSelection(doc.Selection), pkggq.NewNormSelection(doc.Find("h4").First())
}

func TestReturnType_Value(t *testing.T) {
	msg := types.NewNamed("Message", types.KindObject)
	chat := types.NewNamed("Chat", types.KindObject)
	integer := types.NewNamed("Integer", types.KindPrimitive)
	cases := []struct {
		name    string
		desc    string
		noDesc  bool
		want    types.Expression
		wantErr bool
	}{
		{
			name: "returns Union when description matches X is returned otherwise Y is returned",
			desc: "the Message is returned, otherwise Chat is returned",
			want: types.NewUnion(msg, chat),
		},
		{
			name: "returns Array when description contains Array of pattern",
			desc: "Returns an Array of Integer values.",
			want: types.NewArray(integer),
		},
		{
			name: "returns Named when description matches returns as a X object pattern",
			desc: "Returns the result as a Message object.",
			want: msg,
		},
		{
			name: "returns Named when description contains in form of a X object",
			desc: "Returns the message in form of a Message object.",
			want: msg,
		},
		{
			name: "returns Named when description matches a X object pattern",
			desc: "On success, a Message object is returned.",
			want: msg,
		},
		{
			name: "returns Named when description matches Returns X on success",
			desc: "Returns Message on success.",
			want: msg,
		},
		{
			name: "returns Named when description matches Returns the word X pattern",
			desc: "Returns the sent Message.",
			want: msg,
		},
		{
			name: "returns Named when description matches the X is returned without otherwise clause",
			desc: "the Message is returned",
			want: msg,
		},
		{
			name: "returns Named when description matches X is returned fallback pattern",
			desc: "Message is returned.",
			want: msg,
		},
		{
			name:    "returns error when the method has no description paragraphs",
			noDesc:  true,
			wantErr: true,
		},
		{
			name:    "returns error when the return type name is not in the catalog",
			desc:    "Returns UnknownType on success.",
			wantErr: true,
		},
		{
			name:    "returns error when no pattern matches the description text",
			desc:    "Something happened without any recognizable pattern.",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var root, h4 pkggq.Selection
			if tc.noDesc {
				root, h4 = returnTypeNoParaFixture()
			} else {
				root, h4 = returnTypeFixture(tc.desc)
			}
			got, err := gq.NewReturnType(root, h4).Value()
			if tc.wantErr {
				assert.Error(
					t,
					err,
					"ReturnType.Value must return an error when the description cannot be matched to a return type",
				)
				return
			}
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.want,
				got,
				"ReturnType.Value must extract the correct return type expression from the method description",
			)
		})
	}
}
