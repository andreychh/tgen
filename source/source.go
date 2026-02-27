// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package source abstracts data acquisition from different transport layers.
package source

import (
	"context"
	"io"
)

// Source represents a provider of a data stream.
type Source interface {
	// Open initiates access to the resource. The caller must close the stream.
	Open(ctx context.Context) (io.ReadCloser, error)
}
