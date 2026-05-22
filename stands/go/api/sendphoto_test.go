package api_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

func TestSendPhotoMethod_Call(t *testing.T) {
	cases := []struct {
		name   string
		method api.SendPhotoMethod
		check  func(*testing.T, *MultipartCapture)
	}{
		{
			name: "serializes Upload as file part at photo key",
			method: api.SendPhotoMethod{
				ChatID: api.ID(99),
				Photo:  api.Upload{Name: "photo.jpg", Reader: strings.NewReader("data")},
			},
			check: func(t *testing.T, cap *MultipartCapture) {
				part, ok := cap.Part("photo")
				require.True(t, ok, "CapturingConnection must capture the photo field")
				fp, ok := part.(FilePart)
				require.True(t, ok, "Upload must produce a FilePart, not a TextPart")
				assert.Equal(t, "photo.jpg", fp.Name, "FilePart must preserve the Upload filename")
			},
		},
		{
			name: "serializes FileID as text part at photo key",
			method: api.SendPhotoMethod{
				ChatID: api.ID(99),
				Photo:  api.FileID("AgACAgIAAxk"),
			},
			check: func(t *testing.T, cap *MultipartCapture) {
				part, ok := cap.Part("photo")
				require.True(t, ok, "CapturingConnection must capture the photo field")
				tp, ok := part.(TextPart)
				require.True(t, ok, "FileID must produce a TextPart, not a FilePart")
				assert.Equal(t, "AgACAgIAAxk", tp.Value, "TextPart must carry the FileID string verbatim")
			},
		},
		{
			name: "serializes integer ChatID as decimal string",
			method: api.SendPhotoMethod{
				ChatID: api.ID(12345),
				Photo:  api.FileID("x"),
			},
			check: func(t *testing.T, cap *MultipartCapture) {
				part, ok := cap.Part("chat_id")
				require.True(t, ok, "CapturingConnection must capture the chat_id field")
				tp, ok := part.(TextPart)
				require.True(t, ok, "ChatID must produce a TextPart")
				assert.Equal(t, "12345", tp.Value, "integer ChatID must be serialized as a decimal string")
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
