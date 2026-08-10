// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

func TestSendPhotoMethod_Call(t *testing.T) {
	cases := []struct {
		name   string
		chatID api.ChatID
		photo  api.InputFile
		check  func(*testing.T, *http.Request)
	}{
		{
			name:   "sends an upload as a part of its own",
			chatID: api.ID(-1001234567890),
			photo:  api.Upload{Name: "морской ёж.jpg", Reader: strings.NewReader("🐙")},
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				file, ok := requireForm(t, req).File("photo")
				require.True(t, ok, "a parameter that is a file must travel as a part named by the parameter's own key")
				assert.Equal(
					t,
					FormFile{Name: "морской ёж.jpg", Content: "🐙"},
					file,
					"a part must carry the bytes of the upload under the name the upload travels as",
				)
			},
		},
		{
			name:   "leaves nothing behind in the body for an upload",
			chatID: api.ID(-1001234567890),
			photo:  api.Upload{Name: "морской ёж.jpg", Reader: strings.NewReader("🐙")},
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				_, ok := requireForm(t, req).Field("photo")
				assert.False(t, ok, "an upload handed over to a part of its own must leave no value in the body")
			},
		},
		{
			name:   "keeps the parameters that are not files in the body",
			chatID: api.ID(-1001234567890),
			photo:  api.Upload{Name: "морской ёж.jpg", Reader: strings.NewReader("🐙")},
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				value, ok := requireForm(t, req).Field("chat_id")
				require.True(t, ok, "a parameter that is not a file must ride in the body of the same request")
				assert.Equal(t, "-1001234567890", value, "a numeric parameter must ride in the body as the number it is")
			},
		},
		{
			name:   "keeps a file id in the body as a value of its own",
			chatID: api.ID(-1001234567890),
			photo:  api.FileID("AgACAgIAAxkBAAIBY2Vw"),
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.JSONEq(
					t,
					`{"chat_id":-1001234567890,"photo":"AgACAgIAAxkBAAIBY2Vw"}`,
					string(body),
					"a file already held by Telegram must ride in the body as the name it is known by, attaching nothing to the request",
				)
			},
		},
		{
			name:   "addresses a chat by username as the bare string a username is",
			chatID: api.Username("@ёжики"),
			photo:  api.Upload{Name: "морской ёж.jpg", Reader: strings.NewReader("🐙")},
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				value, ok := requireForm(t, req).Field("chat_id")
				require.True(t, ok, "a parameter that is not a file must ride in the body of the same request")
				assert.Equal(t, "@ёжики", value, "a textual parameter must ride in the body unquoted, as the form field it became")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendPhotoMethod{
				ChatID: tc.chatID,
				Photo:  tc.photo,
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, conn.Request())
		})
	}
}

func TestSendDocumentMethod_Call(t *testing.T) {
	cases := []struct {
		name     string
		markup   api.ReplyMarkup
		entities []api.MessageEntity
		key      string
		want     string
	}{
		{
			name:   "carries an object parameter into the body as the JSON it is",
			markup: api.ReplyKeyboardRemove{RemoveKeyboard: true},
			key:    "reply_markup",
			want:   `{"remove_keyboard":true}`,
		},
		{
			name:     "carries a sequence parameter into the body as the JSON it is",
			entities: []api.MessageEntity{{Type: "bold", Offset: 0, Length: 4}},
			key:      "caption_entities",
			want:     `[{"type":"bold","offset":0,"length":4}]`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendDocumentMethod{
				ChatID:          api.ID(-1001234567890),
				Document:        api.Upload{Name: "смета.pdf", Reader: strings.NewReader("%PDF")},
				ReplyMarkup:     tc.markup,
				CaptionEntities: tc.entities,
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			value, ok := requireForm(t, conn.Request()).Field(tc.key)
			require.True(t, ok, "a parameter that is not a file must ride in the body of the same request")
			assert.JSONEq(t, tc.want, value, "a parameter that is no string must ride in the body verbatim, never quoted as one")
		})
	}
}

func TestSendMediaGroupMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		media []api.InputMediaGroup
		check func(*testing.T, Form)
	}{
		{
			name: "attaches a nested upload as the part its reference names",
			media: []api.InputMediaGroup{
				api.InputMediaPhoto{Media: api.Upload{Name: "кот.jpg", Reader: strings.NewReader("JPEG")}},
			},
			check: func(t *testing.T, form Form) {
				t.Helper()
				file, ok := form.File("attachment_0")
				require.True(t, ok, "a file referenced as attach:// must be attached under the key the reference names")
				assert.Equal(
					t,
					FormFile{Name: "кот.jpg", Content: "JPEG"},
					file,
					"a part must carry the bytes of the upload under the name the upload travels as",
				)
			},
		},
		{
			name: "numbers the attachments of one request in turn",
			media: []api.InputMediaGroup{
				api.InputMediaPhoto{Media: api.Upload{Name: "кот.jpg", Reader: strings.NewReader("JPEG")}},
				api.InputMediaPhoto{Media: api.Upload{Name: "пёс.jpg", Reader: strings.NewReader("JFIF")}},
			},
			check: func(t *testing.T, form Form) {
				t.Helper()
				assert.Equal(
					t,
					[]string{"attachment_0", "attachment_1"},
					form.Attached(),
					"every file of one request must reach a key of its own, numbered as the request is walked",
				)
			},
		},
		{
			name: "points every variant at the part its own upload became",
			media: []api.InputMediaGroup{
				api.InputMediaPhoto{Media: api.Upload{Name: "кот.jpg", Reader: strings.NewReader("JPEG")}},
				api.InputMediaPhoto{Media: api.Upload{Name: "пёс.jpg", Reader: strings.NewReader("JFIF")}},
			},
			check: func(t *testing.T, form Form) {
				t.Helper()
				value, ok := form.Field("media")
				require.True(t, ok, "a parameter holding files inside must ride in the body as a value of its own")
				assert.JSONEq(
					t,
					`[{"type":"photo","media":"attach://attachment_0"},{"type":"photo","media":"attach://attachment_1"}]`,
					value,
					"each element of a sequence must reference the part its own file was attached as, keeping the discriminator its marshaller would have written",
				)
			},
		},
		{
			name: "leaves out an optional nested file that was not given",
			media: []api.InputMediaGroup{
				api.InputMediaVideo{Media: api.Upload{Name: "ролик.mp4", Reader: strings.NewReader("ftyp")}},
			},
			check: func(t *testing.T, form Form) {
				t.Helper()
				value, ok := form.Field("media")
				require.True(t, ok, "a parameter holding files inside must ride in the body as a value of its own")
				assert.JSONEq(
					t,
					`[{"type":"video","media":"attach://attachment_0"}]`,
					value,
					"an optional file left unset must leave no reference behind, not even an empty one",
				)
			},
		},
		{
			name: "hands over an optional nested file once it is given",
			media: []api.InputMediaGroup{
				api.InputMediaVideo{
					Media:     api.Upload{Name: "ролик.mp4", Reader: strings.NewReader("ftyp")},
					Thumbnail: api.Upload{Name: "кадр.png", Reader: strings.NewReader("\x89PNG")},
				},
			},
			check: func(t *testing.T, form Form) {
				t.Helper()
				value, ok := form.Field("media")
				require.True(t, ok, "a parameter holding files inside must ride in the body as a value of its own")
				assert.JSONEq(
					t,
					`[{"type":"video","media":"attach://attachment_0","thumbnail":"attach://attachment_1"}]`,
					value,
					"an optional file given inside a variant must reach a part of its own, referenced beside the required one",
				)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendMediaGroupMethod{
				ChatID: api.ID(-1001234567890),
				Media:  tc.media,
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, requireForm(t, conn.Request()))
		})
	}
}

func TestSetMyProfilePhotoMethod_Call(t *testing.T) {
	cases := []struct {
		name  string
		photo api.InputProfilePhoto
		check func(*testing.T, *http.Request)
	}{
		{
			name:  "rewrites a parameter holding an upload into a reference to the part",
			photo: api.InputProfilePhotoStatic{Photo: api.Upload{Name: "аватар.jpg", Reader: strings.NewReader("JPEG")}},
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				value, ok := requireForm(t, req).Field("photo")
				require.True(t, ok, "a parameter holding a file inside must ride in the body as a value of its own")
				assert.JSONEq(
					t,
					`{"type":"static","photo":"attach://attachment_0"}`,
					value,
					"a parameter holding a file must rewrite itself into JSON referencing the part its file became",
				)
			},
		},
		{
			name:  "keeps a file id held inside a parameter as the name it is known by",
			photo: api.InputProfilePhotoStatic{Photo: api.FileID("AgACAgIAAxkBAAIBY2Vw")},
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.JSONEq(
					t,
					`{"photo":{"type":"static","photo":"AgACAgIAAxkBAAIBY2Vw"}}`,
					string(body),
					"a file already held by Telegram must be named where it stands, attaching nothing to the request",
				)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			err := api.SetMyProfilePhotoMethod{Photo: tc.photo}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, conn.Request())
		})
	}
}

func TestSendPollMethod_Call(t *testing.T) {
	cases := []struct {
		name    string
		options []api.InputPollOption
		media   api.InputPollMedia
		check   func(*testing.T, *http.Request)
	}{
		{
			name:    "leaves out an optional parameter that was not given",
			options: []api.InputPollOption{{Text: "да"}, {Text: "нет"}},
			media:   nil,
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.JSONEq(
					t,
					`{"chat_id":-1001234567890,"question":"обед?","options":[{"text":"да"},{"text":"нет"}]}`,
					string(body),
					"an optional parameter that could have held a file must leave the body untouched when it is not given",
				)
			},
		},
		{
			name:    "rewrites an optional parameter holding an upload once it is given",
			options: []api.InputPollOption{{Text: "да"}, {Text: "нет"}},
			media:   api.InputMediaPhoto{Media: api.Upload{Name: "меню.jpg", Reader: strings.NewReader("JPEG")}},
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				value, ok := requireForm(t, req).Field("explanation_media")
				require.True(t, ok, "an optional parameter holding a file must ride in the body as a value of its own once it is given")
				assert.JSONEq(
					t,
					`{"type":"photo","media":"attach://attachment_0"}`,
					value,
					"an optional parameter holding a file must rewrite itself into JSON referencing the part its file became",
				)
			},
		},
		{
			name: "rewrites an element that holds a file, telling it apart by nothing",
			options: []api.InputPollOption{
				{Text: "да", Media: api.InputMediaPhoto{Media: api.Upload{Name: "меню.jpg", Reader: strings.NewReader("JPEG")}}},
				{Text: "нет"},
			},
			media: nil,
			check: func(t *testing.T, req *http.Request) {
				t.Helper()
				value, ok := requireForm(t, req).Field("options")
				require.True(t, ok, "a parameter holding files inside must ride in the body as a value of its own")
				assert.JSONEq(
					t,
					`[{"text":"да","media":{"type":"photo","media":"attach://attachment_0"}},{"text":"нет"}]`,
					value,
					"an object no discriminator tells apart must rewrite itself carrying its own fields alone, the file inside it handed over on the way",
				)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendPollMethod{
				ChatID:           api.ID(-1001234567890),
				Question:         "обед?",
				Options:          tc.options,
				ExplanationMedia: tc.media,
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, conn.Request())
		})
	}
}

func TestAnswerInlineQueryMethod_Call(t *testing.T) {
	cases := []struct {
		name    string
		results []api.InlineQueryResult
		want    string
	}{
		{
			name: "marshals a variant that holds no file, discriminator and all",
			results: []api.InlineQueryResult{
				api.InlineQueryResultGame{ID: "r1", GameShortName: "тетрис"},
			},
			want: `{"inline_query_id":"q1","results":[{"type":"game","id":"r1","game_short_name":"тетрис"}]}`,
		},
		{
			name:    "sends an empty sequence as one, not as nothing",
			results: []api.InlineQueryResult{},
			want:    `{"inline_query_id":"q1","results":[]}`,
		},
		{
			name: "writes the content a result carries as the object that content is",
			results: []api.InlineQueryResult{
				api.InlineQueryResultArticle{
					ID:                  "r1",
					Title:               "ёжик",
					InputMessageContent: api.InputTextMessageContent{MessageText: "в тумане"},
				},
			},
			want: `{"inline_query_id":"q1","results":[{"type":"article","id":"r1","title":"ёжик","input_message_content":{"message_text":"в тумане"}}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			err := api.AnswerInlineQueryMethod{
				InlineQueryID: "q1",
				Results:       tc.results,
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			body, err := io.ReadAll(conn.Request().Body)
			require.NoError(t, err)
			assert.JSONEq(
				t,
				tc.want,
				string(body),
				"a variant a carrier union admits but that holds no file must rewrite itself by marshalling, which is what writes its discriminator",
			)
		})
	}
}

func TestSendAudioMethod_Call(t *testing.T) {
	cases := []struct {
		name      string
		thumbnail api.InputFile
		check     func(*testing.T, Form)
	}{
		{
			name:      "sends an optional upload as a part of its own",
			thumbnail: api.Upload{Name: "обложка.png", Reader: strings.NewReader("\x89PNG")},
			check: func(t *testing.T, form Form) {
				t.Helper()
				file, ok := form.File("thumbnail")
				require.True(t, ok, "an optional parameter that is a file must travel as a part of its own once it is given")
				assert.Equal(
					t,
					FormFile{Name: "обложка.png", Content: "\x89PNG"},
					file,
					"a part must carry the bytes of the upload under the name the upload travels as",
				)
			},
		},
		{
			name:      "attaches no part for an optional file left unset",
			thumbnail: nil,
			check: func(t *testing.T, form Form) {
				t.Helper()
				_, ok := form.File("thumbnail")
				assert.False(t, ok, "an optional parameter left unset must attach no part to the request")
			},
		},
		{
			name:      "leaves no value in the body for an optional file left unset",
			thumbnail: nil,
			check: func(t *testing.T, form Form) {
				t.Helper()
				_, ok := form.Field("thumbnail")
				assert.False(t, ok, "an optional parameter left unset must leave no value in the body, not even an empty one")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := api.SendAudioMethod{
				ChatID:    api.ID(-1001234567890),
				Audio:     api.Upload{Name: "трек.mp3", Reader: strings.NewReader("ID3")},
				Thumbnail: tc.thumbnail,
			}.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, requireForm(t, conn.Request()))
		})
	}
}
