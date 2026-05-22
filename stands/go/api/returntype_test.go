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

const (
	userAlice = `{"id":1,"is_bot":false,"first_name":"Alice"}`
	userBob   = `{"id":2,"is_bot":false,"first_name":"Bob"}`
)

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
			result, err := api.EditMessageTextMethod{Text: "hi"}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, result)
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
