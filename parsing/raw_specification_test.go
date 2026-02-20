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

func TestRawSpecification_Unions(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		limit     int
		wantNames []string
	}{
		{
			name: "extracts only union definitions from a document containing mixed entities",
			html: `
				<h4><a class="anchor" href="#getupdates"></a>getUpdates</h4>
			   	<p>Use this method to receive incoming updates</p>
			   	<table class="table"><thead><tr><th>Parameter</th><th>Type</th><th>Required</th><th>Description</th></tr></thead><tbody><tr><td>offset</td><td>Integer</td><td>Optional</td><td>Identifier of the first update to be returned.</td></tr></tbody></table>
			   	
				<h4><a class="anchor" href="#maybeinaccessiblemessage"></a>MaybeInaccessibleMessage</h4>
			   	<p>This object describes a message...</p>
			   	<ul><li><a href="#message">Message</a></li><li><a href="#inaccessiblemessage">InaccessibleMessage</a></li></ul>
			   	
				<h4><a class="anchor" href="#user"></a>User</h4>
			   	<p>This object represents a Telegram user or bot.</p>
			   	<table class="table"><thead><tr><th>Field</th><th>Type</th><th>Description</th></tr></thead><tbody><tr><td>id</td><td>Integer</td><td>Unique identifier.</td></tr></tbody></table>
			   	
				<h4><a class="anchor" href="#messageorigin"></a>MessageOrigin</h4>
			   	<p>This object describes the origin of a message. It can be one of</p>
			   	<ul><li><a href="#messageoriginuser">MessageOriginUser</a></li><li><a href="#messageoriginhiddenuser">MessageOriginHiddenUser</a></li><li><a href="#messageoriginchat">MessageOriginChat</a></li><li><a href="#messageoriginchannel">MessageOriginChannel</a></li></ul>
			`,
			limit:     -1,
			wantNames: []string{"MaybeInaccessibleMessage", "MessageOrigin"},
		},
		{
			name: "yields an empty sequence when the document contains no union definitions",
			html: `
				<h4><a class="anchor" href="#getupdates"></a>getUpdates</h4>
			   	<p>Use this method to receive incoming updates</p>
			   	<table class="table"><thead><tr><th>Parameter</th><th>Type</th><th>Required</th><th>Description</th></tr></thead><tbody><tr><td>offset</td><td>Integer</td><td>Optional</td><td>Identifier of the first update to be returned.</td></tr></tbody></table>

			    <h4><a class="anchor" href="#user"></a>User</h4>
			    <p>This object represents a Telegram user or bot.</p>
			    <table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>
			`,
			limit:     -1,
			wantNames: []string{},
		},
		{
			name: "stops iterating immediately when the yield function returns false",
			html: `
				<h4><a class="anchor" href="#maybeinaccessiblemessage"></a>MaybeInaccessibleMessage</h4>
			   	<p>This object describes a message...</p>
			   	<ul><li><a href="#message">Message</a></li><li><a href="#inaccessiblemessage">InaccessibleMessage</a></li></ul>
			   	
				<h4><a class="anchor" href="#messageorigin"></a>MessageOrigin</h4>
			   	<p>This object describes the origin of a message. It can be one of</p>
			   	<ul><li><a href="#messageoriginuser">MessageOriginUser</a></li><li><a href="#messageoriginhiddenuser">MessageOriginHiddenUser</a></li><li><a href="#messageoriginchat">MessageOriginChat</a></li><li><a href="#messageoriginchannel">MessageOriginChannel</a></li></ul>
			`,
			limit:     1,
			wantNames: []string{"MaybeInaccessibleMessage"},
		},
		{
			name: "bypasses extraction entirely when the yield function immediately returns false",
			html: `
				<h4><a class="anchor" href="#maybeinaccessiblemessage"></a>MaybeInaccessibleMessage</h4>
			   	<p>This object describes a message...</p>
			   	<ul><li><a href="#message">Message</a></li><li><a href="#inaccessiblemessage">InaccessibleMessage</a></li></ul>
			`,
			limit:     0,
			wantNames: []string{},
		},
		{
			name: "ignores html elements that resemble objects but miss structural requirements",
			html: `
				<h4>Just a regular heading without an anchor</h4>
				<p>Should be ignored.</p>
				<table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>

				<h4><a class="anchor" href="#invalid_name"></a>invalid_name_object</h4>
				<p>Should be ignored due to invalid PascalCase name.</p>
				<table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>
			`,
			limit:     -1,
			wantNames: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			spec := parsing.NewRawSpecification(dom.NewHTMLSelection(doc.Selection))
			gotNames := []string{}
			for union := range spec.Unions() {
				if tt.limit != -1 && len(gotNames) >= tt.limit {
					break
				}
				name, err := union.Name()
				require.NoError(t, err, "extracted union name is invalid or missing")
				gotNames = append(gotNames, name)
			}
			assert.Equal(
				t,
				tt.wantNames,
				gotNames,
				"extracted union sequence does not match expected items",
			)
		})
	}
}

func TestRawSpecification_Objects(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		limit     int
		wantNames []string
	}{
		{
			name: "extracts only object definitions from a document containing mixed entities",
			html: `
				<h4><a class="anchor" href="#game"></a>Game</h4>
				<p>This object represents a game.</p>
				<table class="table"><thead><tr><th>Field</th><th>Type</th><th>Description</th></tr></thead><tbody><tr><td>title</td><td>String</td><td>Title of the game</td></tr></table>
			    
				<h4><a class="anchor" href="#maybeinaccessiblemessage"></a>MaybeInaccessibleMessage</h4>
			   	<p>This object describes a message...</p>
			   	<ul><li><a href="#message">Message</a></li><li><a href="#inaccessiblemessage">InaccessibleMessage</a></li></ul>
			   	
			    <h4><a class="anchor" href="#user"></a>User</h4>
			    <p>This object represents a Telegram user or bot.</p>
			    <table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>
			   	
				<h4><a class="anchor" href="#messageorigin"></a>MessageOrigin</h4>
			   	<p>This object describes the origin of a message. It can be one of</p>
			   	<ul><li><a href="#messageoriginuser">MessageOriginUser</a></li><li><a href="#messageoriginhiddenuser">MessageOriginHiddenUser</a></li><li><a href="#messageoriginchat">MessageOriginChat</a></li><li><a href="#messageoriginchannel">MessageOriginChannel</a></li></ul>
			`,
			limit:     -1,
			wantNames: []string{"Game", "User"},
		},
		// TODO: Fix placeholder object parsing.
		// The parser currently fails to extract objects like CallbackGame because they
		// lack an associated <table>. We need to support objects that only have a
		// description.
		// {
		// 	name: "extracts a placeholder object definition without fields",
		// 	html: `
		// 		<h4><a class="anchor" href="#callbackgame"></a>CallbackGame</h4>
		// 	   	<p>A placeholder, currently holds no information...</p>
		// 	`,
		// 	limit:     -1,
		// 	wantNames: []string{"CallbackGame"},
		// },
		{
			name: "stops iterating immediately when the yield function returns false",
			html: `
				<h4><a class="anchor" href="#game"></a>Game</h4>
				<p>This object represents a game.</p>
				<table class="table"><thead><tr><th>Field</th><th>Type</th><th>Description</th></tr></thead><tbody><tr><td>title</td><td>String</td><td>Title of the game</td></tr></table>
			    
			    <h4><a class="anchor" href="#user"></a>User</h4>
			    <p>This object represents a Telegram user or bot.</p>
			    <table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>
			`,
			limit:     1,
			wantNames: []string{"Game"},
		},
		{
			name: "bypasses extraction entirely when the yield function immediately returns false",
			html: `
			    <h4><a class="anchor" href="#user"></a>User</h4>
			    <p>This object represents a Telegram user or bot.</p>
			    <table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>
			`,
			limit:     0,
			wantNames: []string{},
		},
		{
			name: "ignores html elements that resemble objects but miss structural requirements",
			html: `
				<h4>Just a regular heading without an anchor</h4>
				<p>Should be ignored.</p>
				<table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>

				<h4><a class="anchor" href="#invalid_name"></a>invalid_name_object</h4>
				<p>Should be ignored due to invalid PascalCase name.</p>
				<table class="table"><thead><tr><th>Field</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>Integer</td></tr></tbody></table>
			`,
			limit:     -1,
			wantNames: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			spec := parsing.NewRawSpecification(dom.NewHTMLSelection(doc.Selection))
			gotNames := []string{}
			for obj := range spec.Objects() {
				if tt.limit != -1 && len(gotNames) >= tt.limit {
					break
				}
				name, err := obj.Name()
				require.NoError(t, err, "extracted object name is invalid or missing")
				gotNames = append(gotNames, name)
			}
			assert.Equal(
				t,
				tt.wantNames,
				gotNames,
				"extracted object sequence does not match expected items",
			)
		})
	}
}
