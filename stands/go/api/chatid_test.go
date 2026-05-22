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

func TestSendMessageMethod_Call(t *testing.T) {
	cases := []struct {
		name   string
		chatID api.ChatID
		want   string
	}{
		{
			name:   "serializes ID as decimal string",
			chatID: api.ID(12345),
			want:   "12345",
		},
		{
			name:   "serializes Username verbatim",
			chatID: api.Username("@octocat"),
			want:   "@octocat",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendMessageMethod{
				ChatID: tc.chatID,
				Text:   "hi",
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			part, ok := conn.Capture.Part("chat_id")
			require.True(t, ok, "CapturingConnection must capture the chat_id field")
			tp, ok := part.(TextPart)
			require.True(t, ok, "ChatID must produce a TextPart")
			assert.Equal(t, tc.want, tp.Value, "ChatID must serialize to the expected string representation")
		})
	}
}
