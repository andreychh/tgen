// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"stand/api"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func announceHoshiCup(
	ctx context.Context,
	conn api.Connection,
	chatID api.ChatID,
	poster io.Reader,
	rng *rand.Rand,
) error {
	info, err := api.GetChatMethod{ChatID: chatID}.Call(ctx, conn)
	if err != nil {
		return fmt.Errorf("getting chat info: %w", err)
	}
	msg, err := api.SendPhotoMethod{
		ChatID:  chatID,
		Photo:   api.Upload{Name: "poster.jpg", Reader: poster},
		Caption: ptr("Hoshi Cup IX is open for registration."),
		ReplyMarkup: api.InlineKeyboardMarkup{
			InlineKeyboard: [][]api.InlineKeyboardButton{
				{{Text: "Register", URL: ptr("https://hoshicup.ru/register")}},
			},
		},
	}.Call(ctx, conn)
	if err != nil {
		return fmt.Errorf("sending photo: %w", err)
	}
	var found []api.ReactionType
	for _, r := range info.AvailableReactions {
		if _, ok := r.(api.ReactionTypeEmoji); ok {
			found = append(found, r)
		}
	}
	if len(found) == 0 {
		return nil
	}
	_, err = api.SetMessageReactionMethod{
		ChatID:    api.ID(msg.Chat.ID),
		MessageID: msg.MessageID,
		Reaction:  []api.ReactionType{found[rng.Intn(len(found))]},
	}.Call(ctx, conn)
	if err != nil {
		return fmt.Errorf("setting reaction: %w", err)
	}
	return nil
}

func ptr[T any](v T) *T { return &v }

func TestAnnounceHoshiCup(t *testing.T) {
	queue := api.NewSeqCallQueue(
		api.NewCall(api.MethodGetChat, api.Ok(api.ChatFullInfo{
			AvailableReactions: []api.ReactionType{
				api.ReactionTypeEmoji{Emoji: "🎉"},
				api.ReactionTypeEmoji{Emoji: "👏"},
			},
		})),
		api.NewCall(api.MethodSendPhoto, api.Ok(api.Message{
			MessageID: 101,
			Chat:      api.Chat{ID: -1001234567890, Type: "supergroup"},
		})),
		api.NewCall(api.MethodSetMessageReaction, api.Ok(true)),
	)
	conn := api.NewFakeConnection(queue)
	err := announceHoshiCup(
		context.Background(),
		conn,
		api.Username("@hoshi_cup"),
		strings.NewReader("poster"),
		rand.New(rand.NewSource(42)),
	)
	require.NoError(t, err)
	assert.Len(
		t,
		queue.Calls(),
		3,
		"announceHoshiCup must exhaust all three queued calls in sequence",
	)
}

func requireEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("%s not set", key)
	}
	return v
}

func TestAnnounceHoshiCup_Live(t *testing.T) {
	token := requireEnv(t, "BOT_API_TOKEN")
	chatID, err := strconv.ParseInt(requireEnv(t, "CHAT_ID"), 10, 64)
	require.NoError(t, err)
	path := requireEnv(t, "POSTER_PATH")
	poster, err := os.Open(path)
	require.NoError(t, err)
	defer poster.Close()
	conn := api.NewHTTPConnection(http.DefaultClient, token)
	err = announceHoshiCup(
		context.Background(),
		conn,
		api.ID(chatID),
		poster,
		rand.New(rand.NewSource(0)),
	)
	assert.NoError(t, err)
}
