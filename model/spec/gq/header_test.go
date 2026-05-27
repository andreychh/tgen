// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model/spec/gq"
	pkggq "github.com/andreychh/tgen/pkg/gq"
)

func headerFixture(inner string) (pkggq.Selection, pkggq.Selection) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(
		`<div id="dev_page_content">` + inner + `</div>`,
	))
	return pkggq.NewNormSelection(doc.Selection),
		pkggq.NewNormSelection(doc.Find("div#dev_page_content h4").First())
}

func TestHeader_Kind(t *testing.T) {
	cases := []struct {
		name  string
		inner string
		want  gq.DefinitionKind
	}{
		{
			name:  "returns Unknown when the h4 has no anchor element",
			inner: `<h4>SendMessage</h4><p>Description.</p>`,
			want:  gq.DefinitionKindUnknown,
		},
		{
			name:  "returns Unknown when anchor href contains a hyphen",
			inner: `<h4><a class="anchor" href="#available-types">Available types</a></h4>`,
			want:  gq.DefinitionKindUnknown,
		},
		{
			name: "returns Method when h4 text starts with a lowercase letter",
			inner: `<h4><a class="anchor" href="#sendmessage">sendMessage</a></h4>` +
				`<p>Use this to send messages.</p>`,
			want: gq.DefinitionKindMethod,
		},
		{
			name: "returns Object when h4 is uppercase with no list and no discriminator row",
			inner: `<h4><a class="anchor" href="#message">Message</a></h4>` +
				`<table><tbody>` +
				`<tr><td>message_id</td><td>Integer</td><td>Unique message ID.</td></tr>` +
				`</tbody></table>`,
			want: gq.DefinitionKindObject,
		},
		{
			name: "returns Object when h4 has a discriminator row but is not listed in any union",
			inner: `<h4><a class="anchor" href="#standalone">Standalone</a></h4>` +
				`<table><tbody>` +
				"<tr><td>type</td><td>String</td><td>Always “standalone”</td></tr>" +
				`</tbody></table>`,
			want: gq.DefinitionKindObject,
		},
		{
			name: "returns DiscriminatedObject when h4 has a discriminator row and is listed in a union",
			inner: `<h4><a class="anchor" href="#myvariant">MyVariant</a></h4>` +
				`<table><tbody>` +
				"<tr><td>type</td><td>String</td><td>Always “myvariant”</td></tr>" +
				`</tbody></table>` +
				`<h4><a class="anchor" href="#myunion">MyUnion</a></h4>` +
				`<p>Desc.</p>` +
				`<ul><li><a href="#myvariant">MyVariant</a></li></ul>`,
			want: gq.DefinitionKindDiscriminatedObject,
		},
		{
			name: "returns StructuredUnion when list variants are not defined as h4 elements in the document",
			inner: `<h4><a class="anchor" href="#inputmedia">InputMedia</a></h4>` +
				`<p>Desc.</p>` +
				`<ul>` +
				`<li><a href="#inputmediaphoto">InputMediaPhoto</a></li>` +
				`<li><a href="#inputmediavideo">InputMediaVideo</a></li>` +
				`</ul>`,
			want: gq.DefinitionKindStructuredUnion,
		},
		{
			name: "returns StructuredUnion when list variants exist in the document but have no discriminator rows",
			inner: `<h4><a class="anchor" href="#reactiontype">ReactionType</a></h4>` +
				`<p>Desc.</p>` +
				`<ul><li><a href="#reactiontypeemoji">ReactionTypeEmoji</a></li></ul>` +
				`<h4><a class="anchor" href="#reactiontypeemoji">ReactionTypeEmoji</a></h4>` +
				`<table><tbody>` +
				`<tr><td>type</td><td>String</td><td>Type of the reaction.</td></tr>` +
				`</tbody></table>`,
			want: gq.DefinitionKindStructuredUnion,
		},
		{
			name: "returns DiscriminatedUnion when all list variants have unique discriminator values",
			inner: `<h4><a class="anchor" href="#botcommandscope">BotCommandScope</a></h4>` +
				`<p>Desc.</p>` +
				`<ul>` +
				`<li><a href="#botcommandscopedefault">BotCommandScopeDefault</a></li>` +
				`<li><a href="#botcommandscopechat">BotCommandScopeChat</a></li>` +
				`</ul>` +
				`<h4><a class="anchor" href="#botcommandscopedefault">BotCommandScopeDefault</a></h4>` +
				`<table><tbody>` +
				"<tr><td>type</td><td>String</td><td>Always “default”</td></tr>" +
				`</tbody></table>` +
				`<h4><a class="anchor" href="#botcommandscopechat">BotCommandScopeChat</a></h4>` +
				`<table><tbody>` +
				"<tr><td>type</td><td>String</td><td>Always “chat”</td></tr>" +
				`</tbody></table>`,
			want: gq.DefinitionKindDiscriminatedUnion,
		},
		{
			name: "returns GroupedUnion when all list variants have discriminators with duplicate values",
			inner: `<h4><a class="anchor" href="#myunion">MyUnion</a></h4>` +
				`<p>Desc.</p>` +
				`<ul>` +
				`<li><a href="#variantalpha">VariantAlpha</a></li>` +
				`<li><a href="#variantbeta">VariantBeta</a></li>` +
				`</ul>` +
				`<h4><a class="anchor" href="#variantalpha">VariantAlpha</a></h4>` +
				`<table><tbody>` +
				"<tr><td>type</td><td>String</td><td>Always “group”</td></tr>" +
				`</tbody></table>` +
				`<h4><a class="anchor" href="#variantbeta">VariantBeta</a></h4>` +
				`<table><tbody>` +
				"<tr><td>type</td><td>String</td><td>Always “group”</td></tr>" +
				`</tbody></table>`,
			want: gq.DefinitionKindGroupedUnion,
		},
		{
			name: "returns FallbackUnion when only some list variants have discriminator rows",
			inner: `<h4><a class="anchor" href="#myunion">MyUnion</a></h4>` +
				`<p>Desc.</p>` +
				`<ul>` +
				`<li><a href="#variantone">VariantOne</a></li>` +
				`<li><a href="#varianttwo">VariantTwo</a></li>` +
				`</ul>` +
				`<h4><a class="anchor" href="#variantone">VariantOne</a></h4>` +
				`<table><tbody>` +
				"<tr><td>type</td><td>String</td><td>Always “one”</td></tr>" +
				`</tbody></table>` +
				`<h4><a class="anchor" href="#varianttwo">VariantTwo</a></h4>` +
				`<table><tbody>` +
				`<tr><td>type</td><td>String</td><td>Plain description.</td></tr>` +
				`</tbody></table>`,
			want: gq.DefinitionKindFallbackUnion,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, h4 := headerFixture(tc.inner)
			assert.Equal(
				t,
				tc.want,
				gq.NewHeader(root, h4).Kind(),
				"Header.Kind must classify the h4 into the correct DefinitionKind based on its structure and document context",
			)
		})
	}
}
