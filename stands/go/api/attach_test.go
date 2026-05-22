package api_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

// TestPostStoryMethod_Call covers pattern: method → object with a direct file
// field serialized via attach://.
func TestPostStoryMethod_Call(t *testing.T) {
	cases := []struct {
		name   string
		method api.PostStoryMethod
		check  func(*testing.T, *MultipartCapture)
	}{
		{
			name: "links content photo attach reference to the uploaded file part",
			method: api.PostStoryMethod{
				BusinessConnectionID: "bc1",
				ActivePeriod:         86400,
				Content: api.InputStoryContentPhoto{
					Photo: api.Upload{Name: "story.jpg", Reader: strings.NewReader("data")},
				},
			},
			check: func(t *testing.T, cap *MultipartCapture) {
				raw, ok := cap.Part("content")
				require.True(t, ok, "CapturingConnection must capture the content field")
				tp, ok := raw.(TextPart)
				require.True(t, ok, "content must be serialized as a TextPart")
				var body struct {
					Photo string `json:"photo"`
				}
				require.NoError(t, json.Unmarshal([]byte(tp.Value), &body))
				require.True(t, strings.HasPrefix(body.Photo, "attach://"), "content JSON photo field must be an attach reference")
				key := strings.TrimPrefix(body.Photo, "attach://")
				part, ok := cap.Part(key)
				require.True(t, ok, "file must be registered under the key referenced in the JSON")
				fp, ok := part.(FilePart)
				require.True(t, ok, "the referenced part must be a FilePart, not a TextPart")
				assert.Equal(t, "story.jpg", fp.Name, "FilePart must preserve the Upload filename")
			},
		},
		{
			name: "uses FileID verbatim in content JSON without attach reference",
			method: api.PostStoryMethod{
				BusinessConnectionID: "bc1",
				ActivePeriod:         86400,
				Content: api.InputStoryContentPhoto{
					Photo: api.FileID("AgACAgIAAxk"),
				},
			},
			check: func(t *testing.T, cap *MultipartCapture) {
				raw, ok := cap.Part("content")
				require.True(t, ok, "CapturingConnection must capture the content field")
				tp, ok := raw.(TextPart)
				require.True(t, ok, "content must be serialized as a TextPart")
				var body struct {
					Photo string `json:"photo"`
				}
				require.NoError(t, json.Unmarshal([]byte(tp.Value), &body))
				assert.Equal(t, "AgACAgIAAxk", body.Photo, "FileID must appear verbatim in content JSON without attach:// wrapper")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := tc.method.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, conn.Capture)
		})
	}
}

// TestSendMediaGroupMethod_Call covers pattern: method → union of objects with
// file fields, each serialized via attach://.
func TestSendMediaGroupMethod_Call(t *testing.T) {
	cases := []struct {
		name   string
		method api.SendMediaGroupMethod
		check  func(*testing.T, *MultipartCapture)
	}{
		{
			name: "links each InputMediaPhoto to a distinct file part via attach",
			method: api.SendMediaGroupMethod{
				ChatID: api.ID(1),
				Media: []api.InputMediaGroup{
					api.InputMediaPhoto{Media: api.Upload{Name: "photo1.jpg", Reader: strings.NewReader("d1")}},
					api.InputMediaPhoto{Media: api.Upload{Name: "photo2.jpg", Reader: strings.NewReader("d2")}},
				},
			},
			check: func(t *testing.T, cap *MultipartCapture) {
				raw, ok := cap.Part("media")
				require.True(t, ok, "CapturingConnection must capture the media field")
				tp, ok := raw.(TextPart)
				require.True(t, ok, "media must be serialized as a TextPart")
				var items []struct {
					Media string `json:"media"`
				}
				require.NoError(t, json.Unmarshal([]byte(tp.Value), &items))
				require.Len(t, items, 2, "media array must contain two elements")
				keys := make([]string, len(items))
				for i, item := range items {
					require.True(t, strings.HasPrefix(item.Media, "attach://"), "item %d must be an attach reference", i)
					key := strings.TrimPrefix(item.Media, "attach://")
					part, ok := cap.Part(key)
					require.True(t, ok, "item %d: file must be registered at the referenced key", i)
					_, ok = part.(FilePart)
					require.True(t, ok, "item %d: referenced part must be a FilePart", i)
					keys[i] = key
				}
				assert.NotEqual(t, keys[0], keys[1], "each media item must reference a distinct attach key")
			},
		},
		{
			name: "uses FileID verbatim in media JSON without attach reference",
			method: api.SendMediaGroupMethod{
				ChatID: api.ID(1),
				Media: []api.InputMediaGroup{
					api.InputMediaPhoto{Media: api.FileID("AgACAgIAAxk1")},
					api.InputMediaPhoto{Media: api.FileID("AgACAgIAAxk2")},
				},
			},
			check: func(t *testing.T, cap *MultipartCapture) {
				raw, ok := cap.Part("media")
				require.True(t, ok, "CapturingConnection must capture the media field")
				tp, ok := raw.(TextPart)
				require.True(t, ok, "media must be serialized as a TextPart")
				var items []struct {
					Media string `json:"media"`
				}
				require.NoError(t, json.Unmarshal([]byte(tp.Value), &items))
				require.Len(t, items, 2, "media array must contain two elements")
				assert.Equal(t, "AgACAgIAAxk1", items[0].Media, "FileID must appear verbatim in media JSON element without attach:// wrapper")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := NewCapturingConnection()
			_, err := tc.method.Call(context.Background(), conn)
			require.NoError(t, err)
			tc.check(t, conn.Capture)
		})
	}
}
