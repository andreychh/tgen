// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient abstracts HTTP request execution, primarily for testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPSource retrieves data from a remote URL.
type HTTPSource struct {
	url    string
	client HTTPClient
}

// NewHTTPSource creates an HTTPSource with the given client.
func NewHTTPSource(url string, client HTTPClient) HTTPSource {
	return HTTPSource{url: url, client: client}
}

// NewDefaultHTTPSource creates an HTTPSource using the default [http.Client].
func NewDefaultHTTPSource(url string) HTTPSource {
	return NewHTTPSource(
		url,
		&http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
	)
}

// Open executes a GET request and returns the response body.
func (s HTTPSource) Open(ctx context.Context) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request for %q: %w", s.url, err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request to %q: %w", s.url, err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code %d for %q", resp.StatusCode, s.url)
	}
	return resp.Body, nil
}
