// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

func TestEditMessageTextMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, api.MaybeMessage, error)
	}{
		{
			name: "reads an object back as the message that was edited",
			data: `{"message_id":42,"date":1700000000,"chat":{"id":1,"type":"private"}}`,
			check: func(t *testing.T, result api.MaybeMessage, err error) {
				t.Helper()
				require.NoError(t, err)
				message, ok := result.(api.Message)
				require.True(t, ok, "an object must be read back as Message, the variant an object is")
				assert.Equal(t, int64(42), message.MessageID, "the message read back must carry what the response held")
			},
		},
		{
			name: "reads a boolean back as the marker an inline message answers with",
			data: `true`,
			check: func(t *testing.T, result api.MaybeMessage, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, api.True(true), result, "a boolean must be read back as True, an object never decoding as one")
			},
		},
		{
			name: "refuses a payload that is neither",
			data: `"ни то ни сё"`,
			check: func(t *testing.T, _ api.MaybeMessage, err error) {
				t.Helper()
				assert.ErrorContains(t, err, "cannot unmarshal", "a payload that is neither variant must be refused, not guessed at")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := api.EditMessageTextMethod{}.Call(
				context.Background(),
				NewRespondingConnection(tc.data),
			)
			tc.check(t, result, err)
		})
	}
}

func TestForwardMessageMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, api.Message, error)
	}{
		{
			name: "reads a message dated from zero back as the one beyond reach",
			data: message(`"pinned_message":{"message_id":7,"date":0,"chat":{"id":1,"type":"private"}}`),
			check: func(t *testing.T, message api.Message, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.IsType(t, api.InaccessibleMessage{}, message.PinnedMessage, "a message dated from zero must be read back as the inaccessible one, both variants being objects")
			},
		},
		{
			name: "reads a message dated from anything else back as a message",
			data: message(`"pinned_message":{"message_id":7,"date":1700000000,"chat":{"id":1,"type":"private"}}`),
			check: func(t *testing.T, message api.Message, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.IsType(t, api.Message{}, message.PinnedMessage, "a message dated from anything but zero must be read back as a message")
			},
		},
		{
			name: "reads a bare string back as a plain run of text",
			data: message(`"rich_message":{"blocks":[{"type":"paragraph","text":"простой текст"}]}`),
			check: func(t *testing.T, message api.Message, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, api.RichTextPlain("простой текст"), paragraph(t, message), "a JSON string must be read back as the plain run of text it is")
			},
		},
		{
			name: "reads an array back as a sequence, element by element",
			data: message(`"rich_message":{"blocks":[{"type":"paragraph","text":["начало, ","конец"]}]}`),
			check: func(t *testing.T, message api.Message, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(
					t,
					api.RichTextSequence{api.RichTextPlain("начало, "), api.RichTextPlain("конец")},
					paragraph(t, message),
					"a JSON array must be read back as a sequence, each element read back as itself",
				)
			},
		},
		{
			name: "reads a marked run back as the variant its key names",
			data: message(`"rich_message":{"blocks":[{"type":"paragraph","text":{"type":"bold","text":"жирный"}}]}`),
			check: func(t *testing.T, message api.Message, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(
					t,
					api.RichTextBold{Text: api.RichTextPlain("жирный")},
					paragraph(t, message),
					"a JSON object must be read back as the variant its key names, the run it marks read back in turn",
				)
			},
		},
		{
			name: "reads a sequence back whatever shapes its elements are",
			data: message(`"rich_message":{"blocks":[{"type":"paragraph","text":["начало, ",{"type":"italic","text":"курсив"}]}]}`),
			check: func(t *testing.T, message api.Message, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(
					t,
					api.RichTextSequence{
						api.RichTextPlain("начало, "),
						api.RichTextItalic{Text: api.RichTextPlain("курсив")},
					},
					paragraph(t, message),
					"a sequence must read every element back as the shape that element is, marked runs among plain ones",
				)
			},
		},
		{
			name: "refuses a key no variant answers to",
			data: message(`"rich_message":{"blocks":[{"type":"paragraph","text":{"type":"кот","text":"мяу"}}]}`),
			check: func(t *testing.T, _ api.Message, err error) {
				t.Helper()
				assert.ErrorContains(t, err, "unknown RichText", "a key no variant answers to must be refused, not guessed at")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			message, err := api.ForwardMessageMethod{
				ChatID:     api.ID(1),
				FromChatID: api.ID(2),
				MessageID:  1,
			}.Call(context.Background(), NewRespondingConnection(tc.data))
			tc.check(t, message, err)
		})
	}
}

// message wraps one field into the smallest message a response can hold.
func message(field string) string {
	return `{"message_id":1,"date":1700000000,"chat":{"id":1,"type":"private"},` + field + `}`
}

// paragraph returns the text of the first block of a message, failing the test
// when the message holds no paragraph there.
func paragraph(t *testing.T, message api.Message) api.RichText {
	t.Helper()
	require.NotNil(t, message.RichMessage, "the response must have populated the rich message")
	require.NotEmpty(t, message.RichMessage.Blocks, "the rich message must carry the block the response held")
	block, ok := message.RichMessage.Blocks[0].(api.RichBlockParagraph)
	require.True(t, ok, "the block must have been dispatched to the variant its key names")
	return block.Text
}
