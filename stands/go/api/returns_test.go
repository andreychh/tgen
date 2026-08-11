// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

const (
	userAlisa = `{"id":1,"is_bot":false,"first_name":"Алиса"}`
	userBoris = `{"id":2,"is_bot":false,"first_name":"Борис"}`
)

func TestGetChatMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, api.ChatFullInfo)
	}{
		{
			name: "decodes the value the API answered with",
			data: `{"id":-1001234567890,"type":"supergroup","title":"ежи"}`,
			check: func(t *testing.T, info api.ChatFullInfo) {
				t.Helper()
				assert.Equal(t, int64(-1001234567890), info.ID, "a method returning a value must hand back what the response held")
			},
		},
		{
			name: "dispatches every element of a union field to the variant its key names",
			data: `{"id":1,"type":"private","available_reactions":[{"type":"emoji","emoji":"👍"},{"type":"paid"}]}`,
			check: func(t *testing.T, info api.ChatFullInfo) {
				t.Helper()
				assert.Equal(
					t,
					[]api.ReactionType{api.ReactionTypeEmoji{Emoji: "👍"}, api.ReactionTypePaid{}},
					info.AvailableReactions,
					"a field holding a sequence of a union must decode element by element, each into the variant its key names",
				)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := api.GetChatMethod{ChatID: api.ID(1)}.Call(
				context.Background(),
				NewRespondingConnection(tc.data),
			)
			require.NoError(t, err)
			tc.check(t, info)
		})
	}
}

func TestGetChatMemberMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, api.ChatMember, error)
	}{
		{
			name: "dispatches a payload to the variant its key names",
			data: `{"status":"creator","user":` + userAlisa + `,"is_anonymous":true}`,
			check: func(t *testing.T, member api.ChatMember, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(
					t,
					api.ChatMemberOwner{
						User:        api.User{ID: 1, IsBot: false, FirstName: "Алиса"},
						IsAnonymous: true,
					},
					member,
					"a returned union must decode into the variant its key names, carrying every field of it",
				)
			},
		},
		{
			name: "tells two variants apart by the key alone",
			data: `{"status":"member","user":` + userBoris + `}`,
			check: func(t *testing.T, member api.ChatMember, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.IsType(t, api.ChatMemberMember{}, member, "a key naming another variant must decode into that variant and no other")
			},
		},
		{
			name: "refuses a key no variant answers to",
			data: `{"status":"кот","user":` + userAlisa + `}`,
			check: func(t *testing.T, _ api.ChatMember, err error) {
				t.Helper()
				assert.ErrorContains(t, err, "unknown ChatMember", "a key no variant answers to must be refused, not guessed at")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			member, err := api.GetChatMemberMethod{ChatID: api.ID(1), UserID: 1}.Call(
				context.Background(),
				NewRespondingConnection(tc.data),
			)
			tc.check(t, member, err)
		})
	}
}

func TestGetChatAdministratorsMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		check func(*testing.T, []api.ChatMember, error)
	}{
		{
			name: "decodes every element into the variant its own key names",
			data: `[{"status":"creator","user":` + userAlisa + `,"is_anonymous":false},{"status":"member","user":` + userBoris + `}]`,
			check: func(t *testing.T, members []api.ChatMember, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(
					t,
					[]api.ChatMember{
						api.ChatMemberOwner{User: api.User{ID: 1, FirstName: "Алиса"}},
						api.ChatMemberMember{User: api.User{ID: 2, FirstName: "Борис"}},
					},
					members,
					"a returned sequence of a union must decode element by element, keeping the order the response gave",
				)
			},
		},
		{
			name: "refuses the whole sequence when one element names no variant",
			data: `[{"status":"creator","user":` + userAlisa + `},{"status":"кот","user":` + userBoris + `}]`,
			check: func(t *testing.T, members []api.ChatMember, err error) {
				t.Helper()
				require.ErrorContains(t, err, "unknown ChatMember")
				assert.Nil(t, members, "an element that names no variant must fail the call outright, not leave a half-decoded sequence behind")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			members, err := api.GetChatAdministratorsMethod{ChatID: api.ID(1)}.Call(
				context.Background(),
				NewRespondingConnection(tc.data),
			)
			tc.check(t, members, err)
		})
	}
}

func TestCloseMethod_Call(t *testing.T) {
	failure := errors.New("соединение оборвалось")
	cases := []struct {
		name  string
		conn  api.Connection
		check func(*testing.T, error)
	}{
		{
			name: "confirms a call that succeeded, with nothing to hand back",
			conn: NewRespondingConnection(`true`),
			check: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err, "a method that returns nothing but a confirmation must report success as no error at all")
			},
		},
		{
			name: "hands the failure of the connection back untouched",
			conn: NewFailingConnection(failure),
			check: func(t *testing.T, err error) {
				t.Helper()
				assert.ErrorIs(t, err, failure, "a method must hand back the failure the connection reported, not one of its own")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, api.CloseMethod{}.Call(context.Background(), tc.conn))
		})
	}
}
