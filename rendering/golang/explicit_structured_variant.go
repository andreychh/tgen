// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// ExplicitStructuredVariant represents a named variant of an ImplicitDiscriminatedUnion for Go code generation.
type ExplicitStructuredVariant struct {
	inner explicit.StructuredVariant
}

// NewExplicitStructuredVariant constructs an ExplicitStructuredVariant from an explicit variant.
func NewExplicitStructuredVariant(v explicit.StructuredVariant) ExplicitStructuredVariant {
	return ExplicitStructuredVariant{inner: v}
}

func (e ExplicitStructuredVariant) Name() Name {
	return NewName(e.inner.Name())
}

func (e ExplicitStructuredVariant) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(e.inner.Reference(), e.inner.Description()))
}

func (e ExplicitStructuredVariant) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(e.inner.Fields(), NewField)
}
