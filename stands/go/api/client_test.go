// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func staticTrip(body string) roundTripFunc {
	return func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}
}

type emptyPayload struct{}

func (emptyPayload) Request(ctx context.Context, method, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, http.NoBody)
}

func TestError_Error(t *testing.T) {
	cases := []struct {
		name string
		code int
		desc string
		want string
	}{
		{
			name: "combines numeric code and description in telegram format",
			code: 400,
			desc: "Bad Request: message text is empty",
			want: "telegram 400: Bad Request: message text is empty",
		},
		{
			name: "combines zero code with placeholder description",
			code: 0,
			desc: "<no description>",
			want: "telegram 0: <no description>",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := &api.Error{Code: tc.code, Description: tc.desc}
			assert.Equal(t, tc.want, err.Error(), "Error must format code and description as 'telegram <code>: <desc>'")
		})
	}
}

func TestHTTPConnection_Do(t *testing.T) {
	sentinel := errors.New("refused")
	cases := []struct {
		name  string
		trip  roundTripFunc
		check func(*testing.T, bool, error)
	}{
		{
			name: "decodes result field into response value on ok:true envelope",
			trip: staticTrip(`{"ok":true,"result":true}`),
			check: func(t *testing.T, got bool, err error) {
				require.NoError(t, err)
				assert.True(t, got, "HTTPConnection must unmarshal the result field into the response argument")
			},
		},
		{
			name: "returns Error carrying code and description on ok:false envelope",
			trip: staticTrip(`{"ok":false,"error_code":401,"description":"Unauthorized"}`),
			check: func(t *testing.T, _ bool, err error) {
				var target *api.Error
				require.ErrorAs(t, err, &target, "HTTPConnection must return *api.Error when the envelope signals failure")
				assert.Equal(t, api.Error{Code: 401, Description: "Unauthorized"}, *target, "Error must carry the error_code and description from the response envelope")
			},
		},
		{
			name: "returns Error with zero code and placeholder description when fields are absent",
			trip: staticTrip(`{"ok":false}`),
			check: func(t *testing.T, _ bool, err error) {
				var target *api.Error
				require.ErrorAs(t, err, &target, "HTTPConnection must return *api.Error when the envelope signals failure")
				assert.Equal(t, api.Error{Code: 0, Description: "<no description>"}, *target, "Error must default to zero code and '<no description>' when error_code and description are absent from the envelope")
			},
		},
		{
			name: "wraps transport failure with sending-request context",
			trip: func(*http.Request) (*http.Response, error) {
				return nil, sentinel
			},
			check: func(t *testing.T, _ bool, err error) {
				assert.ErrorContains(t, err, "sending request", "HTTPConnection must wrap transport errors with 'sending request' context")
			},
		},
		{
			name: "preserves transport error in wrapping chain for errors.Is",
			trip: func(*http.Request) (*http.Response, error) {
				return nil, sentinel
			},
			check: func(t *testing.T, _ bool, err error) {
				assert.ErrorIs(t, err, sentinel, "HTTPConnection must use %%w when wrapping transport errors so callers can reach the original cause via errors.Is")
			},
		},
		{
			name: "wraps malformed response body with decoding-envelope context",
			trip: staticTrip(`not json`),
			check: func(t *testing.T, _ bool, err error) {
				assert.ErrorContains(t, err, "decoding envelope", "HTTPConnection must wrap JSON envelope decode failures with 'decoding envelope' context")
			},
		},
		{
			name: "wraps undecodable result with decoding-result context",
			trip: staticTrip(`{"ok":true,"result":"not-a-bool"}`),
			check: func(t *testing.T, _ bool, err error) {
				assert.ErrorContains(t, err, "decoding result", "HTTPConnection must wrap result unmarshal failures with 'decoding result' context")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := api.NewHTTPConnection(&http.Client{Transport: tc.trip}, "test-token")
			var got bool
			err := conn.Do(context.Background(), api.Method("getMe"), emptyPayload{}, &got)
			tc.check(t, got, err)
		})
	}
}

func TestFakeConnection_Do(t *testing.T) {
	cases := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "propagates Err response without wrapping",
			run: func(t *testing.T) {
				q := api.NewSeqCallQueue(api.NewCall(api.Method("getMe"), api.Err(errors.New("bang"))))
				conn := api.NewFakeConnection(q)
				err := conn.Do(context.Background(), api.Method("getMe"), emptyPayload{}, new(bool))
				assert.EqualError(t, err, "bang", "FakeConnection must return the error from an Err response verbatim without wrapping")
			},
		},
		{
			name: "decodes Ok value into response argument via JSON round-trip",
			run: func(t *testing.T) {
				q := api.NewSeqCallQueue(api.NewCall(api.Method("getMe"), api.Ok(true)))
				conn := api.NewFakeConnection(q)
				var got bool
				err := conn.Do(context.Background(), api.Method("getMe"), emptyPayload{}, &got)
				require.NoError(t, err)
				assert.True(t, got, "FakeConnection must decode the Ok value into the response argument")
			},
		},
		{
			name: "panics when queue is exhausted",
			run: func(t *testing.T) {
				conn := api.NewFakeConnection(api.NewSeqCallQueue())
				assert.Panics(t, func() {
					_ = conn.Do(context.Background(), api.Method("getMe"), emptyPayload{}, new(bool))
				}, "FakeConnection must panic when the call queue is exhausted rather than returning a silent error")
			},
		},
		{
			name: "panics when called method does not match the queued call",
			run: func(t *testing.T) {
				q := api.NewSeqCallQueue(api.NewCall(api.Method("getMe"), api.Ok(true)))
				conn := api.NewFakeConnection(q)
				assert.Panics(t, func() {
					_ = conn.Do(context.Background(), api.Method("sendMessage"), emptyPayload{}, new(bool))
				}, "FakeConnection must panic when the called method does not match the next queued call")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}

func TestSeqCallQueue_Next(t *testing.T) {
	q := api.NewSeqCallQueue()
	_, err := q.Next()
	assert.Error(t, err, "Next must return an error when the queue is exhausted")
}

func TestError_As(t *testing.T) {
	conn := api.NewHTTPConnection(
		&http.Client{Transport: staticTrip(`{"ok":false,"error_code":403,"description":"Forbidden"}`)},
		"test-token",
	)
	raw := conn.Do(context.Background(), api.Method("getMe"), emptyPayload{}, new(bool))
	err := fmt.Errorf("outer: %w", raw)
	var target *api.Error
	assert.ErrorAs(t, err, &target, "*api.Error must be reachable via errors.As after a caller wraps the error returned by Do")
}