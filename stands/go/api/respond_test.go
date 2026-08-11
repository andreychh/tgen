// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"encoding/json"

	"stand/api"
)

// RespondingConnection implements [api.Connection] by decoding a fixed document
// into the response of every call, which is what HTTPConnection does with the
// result it unwrapped from the envelope. It sends nothing and asks nothing of
// the payload.
type RespondingConnection struct {
	data string
}

// NewRespondingConnection creates a RespondingConnection answering every call
// with data.
func NewRespondingConnection(data string) RespondingConnection {
	return RespondingConnection{data: data}
}

// Do decodes the fixed document into response. It returns an error when the
// document does not fit what the call expects.
func (c RespondingConnection) Do(_ context.Context, _ api.Method, _ api.Payload, response any) error {
	return json.Unmarshal([]byte(c.data), response)
}

// FailingConnection implements [api.Connection] by failing every call with a
// fixed error, so that what a method does with a failure can be watched.
type FailingConnection struct {
	err error
}

// NewFailingConnection creates a FailingConnection failing every call with err.
func NewFailingConnection(err error) FailingConnection {
	return FailingConnection{err: err}
}

// Do returns the fixed error, leaving the response untouched.
func (c FailingConnection) Do(_ context.Context, _ api.Method, _ api.Payload, _ any) error {
	return c.err
}
