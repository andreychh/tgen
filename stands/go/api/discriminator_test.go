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

func TestSetMessageReactionMethod_Call(t *testing.T) {
	cases := []struct {
		name     string
		reaction api.ReactionType
		wantType string
	}{
		{
			name:     "injects type:emoji discriminator for ReactionTypeEmoji",
			reaction: api.ReactionTypeEmoji{Emoji: "👍"},
			wantType: "emoji",
		},
		{
			name:     "injects type:custom_emoji discriminator for ReactionTypeCustomEmoji",
			reaction: api.ReactionTypeCustomEmoji{CustomEmojiID: "5368324170671202286"},
			wantType: "custom_emoji",
		},
		{
			name:     "injects type:paid discriminator for ReactionTypePaid",
			reaction: api.ReactionTypePaid{},
			wantType: "paid",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SetMessageReactionMethod{
				ChatID:    api.ID(99),
				MessageID: 1,
				Reaction:  []api.ReactionType{tc.reaction},
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			part, ok := conn.Capture.Part("reaction")
			require.True(t, ok, "CapturingConnection must capture the reaction field")
			tp, ok := part.(TextPart)
			require.True(t, ok, "reaction must produce a TextPart")
			var items []map[string]any
			require.NoError(t, json.Unmarshal([]byte(tp.Value), &items))
			assert.Equal(t, tc.wantType, items[0]["type"], "ReactionType must inject its type discriminator into serialized JSON without user participation")
		})
	}
}

func TestGetChatMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, api.ChatFullInfo)
	}{
		{
			name: "dispatches type:emoji discriminator to ReactionTypeEmoji",
			data: `{"available_reactions":[{"type":"emoji","emoji":"👍"}]}`,
			check: func(t *testing.T, result api.ChatFullInfo) {
				require.NotEmpty(t, result.AvailableReactions, "response must populate AvailableReactions")
				reaction, ok := result.AvailableReactions[0].(api.ReactionTypeEmoji)
				require.True(t, ok, "type:emoji must dispatch to ReactionTypeEmoji, not another variant")
				assert.Equal(t, "👍", reaction.Emoji, "ReactionTypeEmoji must carry the emoji value from JSON")
			},
		},
		{
			name: "dispatches type:custom_emoji discriminator to ReactionTypeCustomEmoji",
			data: `{"available_reactions":[{"type":"custom_emoji","custom_emoji_id":"5368324170671202286"}]}`,
			check: func(t *testing.T, result api.ChatFullInfo) {
				require.NotEmpty(t, result.AvailableReactions, "response must populate AvailableReactions")
				reaction, ok := result.AvailableReactions[0].(api.ReactionTypeCustomEmoji)
				require.True(t, ok, "type:custom_emoji must dispatch to ReactionTypeCustomEmoji, not another variant")
				assert.Equal(t, "5368324170671202286", reaction.CustomEmojiID, "ReactionTypeCustomEmoji must carry the custom emoji ID from JSON")
			},
		},
		{
			name: "dispatches type:paid discriminator to ReactionTypePaid",
			data: `{"available_reactions":[{"type":"paid"}]}`,
			check: func(t *testing.T, result api.ChatFullInfo) {
				require.NotEmpty(t, result.AvailableReactions, "response must populate AvailableReactions")
				_, ok := result.AvailableReactions[0].(api.ReactionTypePaid)
				assert.True(t, ok, "type:paid must dispatch to ReactionTypePaid, not another variant")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewRespondingConnection([]byte(tc.data))
			result, err := api.GetChatMethod{ChatID: api.ID(99)}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, result)
		})
	}
}
