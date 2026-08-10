// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

func TestLogOutMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		check func(*testing.T, *http.Request)
	}{
		{
			name: "sends a request carrying no body",
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.Empty(t, body, "a method taking no parameter must send nothing as its body")
			},
		},
		{
			name: "sends a request declaring no content type",
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				assert.Empty(
					t,
					req.Header.Get("Content-Type"),
					"a method taking no parameter must declare no content type, having no body to describe",
				)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			err := api.LogOutMethod{}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, conn.Request())
		})
	}
}

func TestSendMessageMethod_Call(t *testing.T) {
	cases := []struct {
		name   string
		chatID api.ChatID
		check  func(*testing.T, *http.Request)
	}{
		{
			name:   "sends a request declaring a JSON content type",
			chatID: api.ID(-1001234567890),
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				assert.Equal(
					t,
					"application/json",
					req.Header.Get("Content-Type"),
					"a method reaching no file must send its parameters as JSON",
				)
			},
		},
		{
			name:   "sends the parameters it was given and no others",
			chatID: api.ID(-1001234567890),
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.JSONEq(
					t,
					`{"chat_id":-1001234567890,"text":"привет, 🐙"}`,
					string(body),
					"a method must marshal itself whole, leaving every optional parameter left unset out of the body",
				)
			},
		},
		{
			name:   "addresses a chat by username as the string a username is",
			chatID: api.Username("@ёжики"),
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.JSONEq(
					t,
					`{"chat_id":"@ёжики","text":"привет, 🐙"}`,
					string(body),
					"either way of addressing a chat must travel as the JSON it is, a number or a string, and never as a union of the two",
				)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendMessageMethod{
				ChatID: tc.chatID,
				Text:   "привет, 🐙",
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, conn.Request())
		})
	}
}

func TestSetMessageReactionMethod_Call(t *testing.T) {
	cases := []struct {
		name     string
		reaction api.ReactionType
		want     string
	}{
		{
			name:     "writes the discriminator of the variant it was given",
			reaction: api.ReactionTypeEmoji{Emoji: "👍"},
			want:     `{"chat_id":-1001234567890,"message_id":42,"reaction":[{"type":"emoji","emoji":"👍"}]}`,
		},
		{
			name:     "writes the discriminator of a variant that carries nothing else",
			reaction: api.ReactionTypePaid{},
			want:     `{"chat_id":-1001234567890,"message_id":42,"reaction":[{"type":"paid"}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			err := api.SetMessageReactionMethod{
				ChatID:    api.ID(-1001234567890),
				MessageID: 42,
				Reaction:  []api.ReactionType{tc.reaction},
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			body, err := io.ReadAll(conn.Request().Body)
			require.NoError(t, err)
			assert.JSONEq(
				t,
				tc.want,
				string(body),
				"a variant a union tells apart must write its own discriminator without the caller naming it",
			)
		})
	}
}

func TestSendRichMessageMethod_Call(t *testing.T) {
	cases := []struct {
		name string
		text api.RichText
		want string
	}{
		{
			name: "writes a plain run of text as the JSON string it is",
			text: api.RichTextPlain("простой текст"),
			want: `"простой текст"`,
		},
		{
			name: "writes a sequence of runs as a JSON array, each element as itself",
			text: api.RichTextSequence{
				api.RichTextPlain("начало, "),
				api.RichTextBold{Text: api.RichTextPlain("жирный")},
			},
			want: `["начало, ",{"type":"bold","text":"жирный"}]`,
		},
		{
			name: "writes a marked run as an object carrying its discriminator",
			text: api.RichTextItalic{Text: api.RichTextPlain("курсив")},
			want: `{"type":"italic","text":"курсив"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendRichMessageMethod{
				ChatID: api.ID(-1001234567890),
				RichMessage: api.InputRichMessage{
					Blocks: []api.InputRichBlock{api.InputRichBlockParagraph{Text: tc.text}},
				},
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			body, err := io.ReadAll(conn.Request().Body)
			require.NoError(t, err)
			assert.JSONEq(
				t,
				`{"chat_id":-1001234567890,"rich_message":{"blocks":[{"type":"paragraph","text":`+tc.want+`}]}}`,
				string(body),
				"every variant of the one union written partly as prose must travel as the JSON shape that variant is",
			)
		})
	}
}
