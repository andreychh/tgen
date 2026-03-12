// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering

import "github.com/andreychh/tgen/meta"

// Meta represents the binary metadata, adapted for use in templates.
type Meta struct {
	inner meta.Meta
}

// NewMeta creates a Meta from the given [meta.Meta].
func NewMeta(m meta.Meta) Meta {
	return Meta{inner: m}
}

// Release returns the release metadata.
func (m Meta) Release() Release {
	return NewRelease(m.inner.Release())
}
