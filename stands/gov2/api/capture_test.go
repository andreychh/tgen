// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"fmt"
	"net/http"

	"stand/api"
)

// CapturingConnection implements [api.Connection] by materializing a payload
// into the request it would have been sent as, and keeping that request instead
// of sending it. It reports no failure and leaves the response untouched, so a
// call through it observes nothing but what it asked to be sent.
type CapturingConnection struct {
	request *http.Request
}

// NewCapturingConnection creates a CapturingConnection holding no request yet.
func NewCapturingConnection() *CapturingConnection {
	return &CapturingConnection{request: nil}
}

// Do materializes payload into a request and keeps it. It returns an error when
// the payload cannot become one.
func (c *CapturingConnection) Do(ctx context.Context, _ api.Method, payload api.Payload, _ any) error {
	req, err := payload.Request(ctx, http.MethodPost, "http://capture")
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	c.request = req
	return nil
}

// Request returns the request the last call was materialized into, or nil when
// no call has been made.
func (c *CapturingConnection) Request() *http.Request {
	return c.request
}
