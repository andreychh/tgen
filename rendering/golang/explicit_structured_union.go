// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// ExplicitStructuredUnion represents a structured union from the Telegram Bot API spec for Go code generation.
type ExplicitStructuredUnion struct {
	inner explicit.StructuredUnion
}

// NewExplicitStructuredUnion constructs an ExplicitStructuredUnion from a parsed structured union.
func NewExplicitStructuredUnion(u explicit.StructuredUnion) ExplicitStructuredUnion {
	return ExplicitStructuredUnion{inner: u}
}

func (u ExplicitStructuredUnion) Name() Name {
	return NewName(u.inner.Name())
}

func (u ExplicitStructuredUnion) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(u.inner.Reference(), u.inner.Description()))
}

func (u ExplicitStructuredUnion) Variants() iter.Seq[StructuredVariant] {
	return iters.NewMappedSeq(
		u.inner.Variants(),
		func(v explicit.StructuredVariant) StructuredVariant {
			return NewExplicitStructuredVariant(v)
		},
	)
}
