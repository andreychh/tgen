package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"stand/api"
)

// CapturedPart is a multipart form field — either a plain-text value or a file
// attachment.
//
//sumtype:decl
type CapturedPart interface{ sealedCapturedPart() }

// TextPart is a plain-text form field.
type TextPart struct{ Value string }

func (TextPart) sealedCapturedPart() {}

// FilePart is a file attachment form field.
type FilePart struct{ Name string }

func (FilePart) sealedCapturedPart()  {}

// MultipartCapture holds all fields of a captured multipart request.
type MultipartCapture struct {
	parts map[string]CapturedPart
}

// Part returns the captured field for key, false if absent.
func (c *MultipartCapture) Part(key string) (CapturedPart, bool) {
	p, ok := c.parts[key]
	return p, ok
}

// CapturingConnection implements api.Connection by materializing the Payload
// into an HTTP request and parsing its multipart body into Capture. It always
// succeeds and leaves the response at its zero value.
type CapturingConnection struct {
	Capture *MultipartCapture
}

// NewCapturingConnection creates a CapturingConnection that records every
// multipart field in its Capture.
func NewCapturingConnection() *CapturingConnection {
	return &CapturingConnection{
		Capture: &MultipartCapture{parts: make(map[string]CapturedPart)},
	}
}

// Do materializes payload and stores its multipart fields in Capture.
func (c *CapturingConnection) Do(ctx context.Context, _ api.Method, payload api.Payload, _ any) error {
	req, err := payload.Request(ctx, http.MethodPost, "http://fake")
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	if ct := req.Header.Get("Content-Type"); ct != "" {
		_, params, err := mime.ParseMediaType(ct)
		if err != nil {
			return fmt.Errorf("parsing content type: %w", err)
		}
		mr := multipart.NewReader(req.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("reading part: %w", err)
			}
			key := p.FormName()
			if fn := p.FileName(); fn != "" {
				c.Capture.parts[key] = FilePart{Name: fn}
			} else {
				data, err := io.ReadAll(p)
				if err != nil {
					return fmt.Errorf("reading field %q: %w", key, err)
				}
				c.Capture.parts[key] = TextPart{Value: string(data)}
			}
		}
	}
	return nil
}

// RespondingConnection implements api.Connection by injecting a fixed JSON
// payload into the response parameter via json.Unmarshal. It mirrors what the
// real HTTPConnection does after unwrapping the Telegram API envelope.
type RespondingConnection struct {
	data []byte
}

// NewRespondingConnection creates a RespondingConnection that unmarshals data
// into the response on every Do call.
func NewRespondingConnection(data []byte) *RespondingConnection {
	return &RespondingConnection{data: data}
}

// Do unmarshals the fixed JSON payload into response.
func (c *RespondingConnection) Do(_ context.Context, _ api.Method, _ api.Payload, response any) error {
	return json.Unmarshal(c.data, response)
}
