// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package source

import (
	"context"
	"io"
	"strings"
)

// LocationSource infers the transport protocol from the location string and
// delegates the Open call to the corresponding implementation.
type LocationSource struct {
	location string
}

// NewLocationSource creates a LocationSource for the provided resource
// identifier.
func NewLocationSource(location string) LocationSource {
	return LocationSource{location: location}
}

// Open resolves the resource and initiates data retrieval based on the format
// of the location string.
func (s LocationSource) Open(ctx context.Context) (io.ReadCloser, error) {
	if s.isURL() {
		return NewDefaultHTTPSource(s.location).Open(ctx)
	}
	return NewFileSource(s.location).Open(ctx)
}

// isURL reports whether the location is a network resource.
func (s LocationSource) isURL() bool {
	return strings.HasPrefix(s.location, "http://") || strings.HasPrefix(s.location, "https://")
}
