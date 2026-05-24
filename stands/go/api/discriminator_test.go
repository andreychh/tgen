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

const (
	userAlice = `{"id":1,"is_bot":false,"first_name":"Alice"}`
	userBob   = `{"id":2,"is_bot":false,"first_name":"Bob"}`
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

func TestGetChatMemberMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, api.ChatMember)
	}{
		{
			name:  "dispatches status:creator to ChatMemberOwner",
			data:  `{"status":"creator","user":` + userAlice + `}`,
			check: func(t *testing.T, result api.ChatMember) {
				member, ok := result.(api.ChatMemberOwner)
				require.True(t, ok, "status:creator must be dispatched to ChatMemberOwner, not another variant")
				assert.Equal(t, int64(1), member.User.ID, "ChatMemberOwner must carry the user ID from the response")
			},
		},
		{
			name:  "dispatches status:member to ChatMemberMember",
			data:  `{"status":"member","user":` + userBob + `}`,
			check: func(t *testing.T, result api.ChatMember) {
				_, ok := result.(api.ChatMemberMember)
				assert.True(t, ok, "status:member must be dispatched to ChatMemberMember, not another variant")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewRespondingConnection([]byte(tc.data))
			result, err := api.GetChatMemberMethod{ChatID: api.ID(99), UserID: 1}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, result)
		})
	}
}

func TestGetChatAdministratorsMethod_Call(t *testing.T) {
	data := `[{"status":"creator","user":` + userAlice + `},{"status":"member","user":` + userBob + `}]`
	cases := []struct {
		name  string
		index int
		check func(*testing.T, api.ChatMember)
	}{
		{
			name:  "dispatches first element as ChatMemberOwner",
			index: 0,
			check: func(t *testing.T, member api.ChatMember) {
				_, ok := member.(api.ChatMemberOwner)
				assert.True(t, ok, "status:creator element must be dispatched to ChatMemberOwner")
			},
		},
		{
			name:  "dispatches second element as ChatMemberMember",
			index: 1,
			check: func(t *testing.T, member api.ChatMember) {
				_, ok := member.(api.ChatMemberMember)
				assert.True(t, ok, "status:member element must be dispatched to ChatMemberMember")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewRespondingConnection([]byte(data))
			result, err := api.GetChatAdministratorsMethod{ChatID: api.ID(99)}.Call(context.Background(), conn)
			require.NoError(t, err)
			require.Len(t, result, 2, "response must contain exactly two members")
			tc.check(t, result[tc.index])
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
