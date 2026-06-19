// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

func TestInlineQueryResult_MarshalJSON(t *testing.T) {
	cases := []struct {
		name      string
		variant   api.InlineQueryResult
		wantType  string
		wantField string
	}{
		{
			name:      "InlineQueryResultPhoto emits type:photo with photo_url",
			variant:   api.InlineQueryResultPhoto{ID: "r1", PhotoURL: "https://example.com/p.jpg", ThumbnailURL: "https://example.com/t.jpg"},
			wantType:  "photo",
			wantField: "photo_url",
		},
		{
			name:      "InlineQueryResultCachedPhoto emits type:photo with photo_file_id",
			variant:   api.InlineQueryResultCachedPhoto{ID: "r1", PhotoFileID: "AgACfile123"},
			wantType:  "photo",
			wantField: "photo_file_id",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.AnswerInlineQueryMethod{
				InlineQueryID: "q1",
				Results:       []api.InlineQueryResult{tc.variant},
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			raw, ok := conn.Capture.Part("results")
			require.True(t, ok, "CapturingConnection must capture the results field")
			tp, ok := raw.(TextPart)
			require.True(t, ok, "results must be serialized as a TextPart")
			var items []map[string]any
			require.NoError(t, json.Unmarshal([]byte(tp.Value), &items))
			require.Len(t, items, 1, "results must contain exactly one element")
			require.Equal(t, tc.wantType, items[0]["type"], "variant must inject the type discriminator")
			_, ok = items[0][tc.wantField]
			assert.True(t, ok, "variant must include its structural field in JSON to allow server-side dispatch within the group")
		})
	}
}

func TestInputMessageContent_MarshalJSON(t *testing.T) {
	cases := []struct {
		name    string
		content api.InputMessageContent
		check   func(*testing.T, map[string]any)
	}{
		{
			name:    "InputTextMessageContent produces JSON with message_text field",
			content: api.InputTextMessageContent{MessageText: "hello"},
			check: func(t *testing.T, body map[string]any) {
				_, ok := body["message_text"]
				assert.True(t, ok, "InputTextMessageContent must produce JSON with a message_text field")
			},
		},
		{
			name:    "InputLocationMessageContent produces JSON with latitude field",
			content: api.InputLocationMessageContent{Latitude: 55.7558, Longitude: 37.6173},
			check: func(t *testing.T, body map[string]any) {
				_, ok := body["latitude"]
				assert.True(t, ok, "InputLocationMessageContent must produce JSON with a latitude field")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.AnswerInlineQueryMethod{
				InlineQueryID: "q1",
				Results: []api.InlineQueryResult{
					api.InlineQueryResultArticle{
						ID:                  "r1",
						Title:               "title",
						InputMessageContent: tc.content,
					},
				},
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			raw, ok := conn.Capture.Part("results")
			require.True(t, ok, "CapturingConnection must capture the results field")
			tp, ok := raw.(TextPart)
			require.True(t, ok, "results must be serialized as a TextPart")
			var items []struct {
				Content map[string]any `json:"input_message_content"`
			}
			require.NoError(t, json.Unmarshal([]byte(tp.Value), &items))
			require.Len(t, items, 1, "results must contain exactly one element")
			tc.check(t, items[0].Content)
		})
	}
}

func TestRichText_Marshal(t *testing.T) {
	cases := []struct {
		name string
		text api.RichText
		want string
	}{
		{
			name: "encodes RichTextPlain as a bare JSON string",
			text: api.RichTextPlain("plain"),
			want: `"plain"`,
		},
		{
			name: "encodes RichTextSequence as a JSON array",
			text: api.RichTextSequence{
				api.RichTextPlain("lead "),
				api.RichTextBold{Text: api.RichTextPlain("bold")},
			},
			want: `["lead ", {"type": "bold", "text": "bold"}]`,
		},
		{
			name: "encodes a named variant as a JSON object carrying its type discriminator",
			text: api.RichTextItalic{Text: api.RichTextPlain("x")},
			want: `{"type": "italic", "text": "x"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(api.RichBlockCaption{Text: tc.text})
			require.NoError(t, err)
			var got struct {
				Text json.RawMessage `json:"text"`
			}
			require.NoError(t, json.Unmarshal(data, &got))
			assert.JSONEq(t, tc.want, string(got.Text), "RichText must serialize each variant to its native JSON shape")
		})
	}
}

func TestRichText_Unmarshal(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		check func(*testing.T, api.RichText)
	}{
		{
			name: "decodes a JSON string into RichTextPlain",
			text: `"plain"`,
			check: func(t *testing.T, v api.RichText) {
				assert.Equal(t, api.RichTextPlain("plain"), v, "a JSON string must deserialize into RichTextPlain")
			},
		},
		{
			name: "decodes a JSON array into RichTextSequence preserving element types",
			text: `["lead ", {"type":"bold","text":"bold"}]`,
			check: func(t *testing.T, v api.RichText) {
				seq, ok := v.(api.RichTextSequence)
				require.True(t, ok, "a JSON array must deserialize into RichTextSequence")
				require.Len(t, seq, 2, "RichTextSequence must preserve every element of the array")
				_, ok = seq[1].(api.RichTextBold)
				assert.True(t, ok, "nested objects inside a RichTextSequence must dispatch by their type discriminator")
			},
		},
		{
			name: "decodes a JSON object into the variant named by its type discriminator",
			text: `{"type":"italic","text":"x"}`,
			check: func(t *testing.T, v api.RichText) {
				_, ok := v.(api.RichTextItalic)
				assert.True(t, ok, "a JSON object must dispatch to the RichText variant named by its type discriminator")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var caption api.RichBlockCaption
			require.NoError(t, json.Unmarshal([]byte(`{"text":`+tc.text+`}`), &caption))
			tc.check(t, caption.Text)
		})
	}
}

func TestMaybeInaccessibleMessage_Unmarshal(t *testing.T) {
	cases := []struct {
		name   string
		pinned string
		check  func(*testing.T, api.MaybeInaccessibleMessage)
	}{
		{
			name:   "dispatches date:0 to InaccessibleMessage",
			pinned: `{"message_id":0,"date":0,"chat":{"id":1,"type":"private"}}`,
			check: func(t *testing.T, msg api.MaybeInaccessibleMessage) {
				_, ok := msg.(api.InaccessibleMessage)
				assert.True(t, ok, "date:0 must be dispatched to InaccessibleMessage, not Message")
			},
		},
		{
			name:   "dispatches non-zero date to Message",
			pinned: `{"message_id":5,"date":1234,"chat":{"id":1,"type":"private"}}`,
			check: func(t *testing.T, msg api.MaybeInaccessibleMessage) {
				m, ok := msg.(api.Message)
				require.True(t, ok, "non-zero date must be dispatched to Message, not InaccessibleMessage")
				assert.Equal(t, int64(5), m.MessageID, "Message must carry the message_id from the response")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := `{"message_id":1,"date":1000,"chat":{"id":1,"type":"private"},"pinned_message":` + tc.pinned + `}`
			conn := NewRespondingConnection([]byte(data))
			result, err := api.SendMessageMethod{ChatID: api.ID(1), Text: "hi"}.Call(context.Background(), conn)
			require.NoError(t, err)
			require.NotNil(t, result.PinnedMessage, "response must populate PinnedMessage")
			tc.check(t, result.PinnedMessage)
		})
	}
}

func TestEditMessageTextMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, api.MaybeMessage)
	}{
		{
			name: "returns Message when API responds with a message object",
			data: `{"message_id":42,"date":0,"chat":{"id":1,"type":"private"}}`,
			check: func(t *testing.T, result api.MaybeMessage) {
				msg, ok := result.(api.Message)
				require.True(t, ok, "object response must be deserialized as Message, not True")
				assert.Equal(t, int64(42), msg.MessageID, "Message must carry the message_id from the response")
			},
		},
		{
			name: "returns True when API responds with boolean true",
			data: `true`,
			check: func(t *testing.T, result api.MaybeMessage) {
				_, ok := result.(api.True)
				assert.True(t, ok, "boolean true response must be deserialized as True, not Message")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewRespondingConnection([]byte(tc.data))
			result, err := api.EditMessageTextMethod{Text: ptr("hi")}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, result)
		})
	}
}
