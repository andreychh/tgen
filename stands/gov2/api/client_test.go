// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stand/api"
)

func TestHTTPConnection_Do(t *testing.T) {
	cases := []struct {
		name        string
		destination func(base, token string) api.Destination
		body        string
		closed      bool
		check       func(*testing.T, string, api.User, error)
	}{
		{
			name:        "posts to the path the token and the method name make",
			destination: api.NewDestination,
			body:        `{"ok":true,"result":{"id":7,"is_bot":true,"first_name":"бот"}}`,
			check: func(t *testing.T, path string, _ api.User, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, "/bot42:XYZ/getMe", path, "a request must reach the endpoint the token and the method name address")
			},
		},
		{
			name:        "posts under a segment of its own in the test environment",
			destination: api.NewTestDestination,
			body:        `{"ok":true,"result":{"id":7,"is_bot":true,"first_name":"бот"}}`,
			check: func(t *testing.T, path string, _ api.User, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, "/bot42:XYZ/test/getMe", path, "the test environment must be addressed by an extra segment after the token")
			},
		},
		{
			name:        "hands back what the envelope carried as its result",
			destination: api.NewDestination,
			body:        `{"ok":true,"result":{"id":7,"is_bot":true,"first_name":"бот"}}`,
			check: func(t *testing.T, _ string, user api.User, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, api.User{ID: 7, IsBot: true, FirstName: "бот"}, user, "the result of an envelope reporting success must be decoded into what the method returns")
			},
		},
		{
			name:        "reports a refusal as the failure the API described",
			destination: api.NewDestination,
			body:        `{"ok":false,"error_code":429,"description":"Too Many Requests","parameters":{"retry_after":30}}`,
			check: func(t *testing.T, _ string, _ api.User, err error) {
				t.Helper()
				var failure *api.Error
				require.ErrorAs(t, err, &failure)
				require.NotNil(t, failure.Parameters, "a refusal carrying parameters must keep them")
				assert.Equal(t, int64(30), *failure.Parameters.RetryAfter, "a refusal must carry back what the API said about retrying")
			},
		},
		{
			name:        "fills in what a refusal left unsaid",
			destination: api.NewDestination,
			body:        `{"ok":false}`,
			check: func(t *testing.T, _ string, _ api.User, err error) {
				t.Helper()
				assert.EqualError(t, err, "telegram 0: <no description>", "a refusal saying nothing must still read as a refusal, not as an empty one")
			},
		},
		{
			name:        "refuses a body that is no envelope",
			destination: api.NewDestination,
			body:        `не json`,
			check: func(t *testing.T, _ string, _ api.User, err error) {
				t.Helper()
				assert.ErrorContains(t, err, "decoding envelope", "a body that is no envelope must fail as one, naming what could not be read")
			},
		},
		{
			name:        "reports a request that never reached the API",
			destination: api.NewDestination,
			body:        `{"ok":true,"result":{}}`,
			closed:      true,
			check: func(t *testing.T, _ string, _ api.User, err error) {
				t.Helper()
				assert.ErrorContains(t, err, "sending request", "a request that never left must fail as a request, not as a refusal of the API")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := ""
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			if tc.closed {
				server.Close()
			}
			conn := api.NewHTTPConnectionTo(server.Client(), tc.destination(server.URL, "42:XYZ"))
			user, err := api.GetMeMethod{}.Call(context.Background(), conn)
			tc.check(t, path, user, err)
		})
	}
}
