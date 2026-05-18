// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedUnion represents a Telegram Bot API discriminated union narrowed
// for code generation.
type DiscriminatedUnion struct {
	inner spec.DiscriminatedUnion
}

// NewDiscriminatedUnion constructs a DiscriminatedUnion from a parsed union.
func NewDiscriminatedUnion(d spec.DiscriminatedUnion) DiscriminatedUnion {
	return DiscriminatedUnion{inner: d}
}

func (d DiscriminatedUnion) Reference() (model.Reference, error) {
	return d.inner.Reference()
}

func (d DiscriminatedUnion) Name() (model.Name, error) {
	return d.inner.Name()
}

func (d DiscriminatedUnion) Description() model.Description {
	return d.inner.Description()
}

func (d DiscriminatedUnion) DiscriminatorKey() (model.Key, error) {
	return d.inner.DiscriminatorKey()
}

func (d DiscriminatedUnion) Variants() iter.Seq[DiscriminatedObject] {
	return iters.NewMappedSeq(d.inner.Variants(), NewDiscriminatedObject)
}
