// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

type StructuredVariant struct {
	li gq.Selection
}

func NewStructuredVariant(li gq.Selection) StructuredVariant {
	return StructuredVariant{li: li}
}

func (v StructuredVariant) Reference() explicit.Reference {
	panic("not implemented")
}

func (v StructuredVariant) Name() explicit.Name {
	panic("not implemented")
}

func (v StructuredVariant) Description() explicit.Description {
	panic("not implemented")
}

func (v StructuredVariant) Fields() iter.Seq[explicit.Field] {
	panic("not implemented")
}
