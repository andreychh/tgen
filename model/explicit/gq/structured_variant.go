// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

type StructuredVariant struct {
	li gq.Selection
}

func NewStructuredVariant(li gq.Selection) StructuredVariant {
	return StructuredVariant{li: li}
}

func (v StructuredVariant) Reference() model.Reference {
	panic("not implemented")
}

func (v StructuredVariant) Name() model.Name {
	panic("not implemented")
}

func (v StructuredVariant) Description() model.Description {
	panic("not implemented")
}

func (v StructuredVariant) Fields() iter.Seq[explicit.Field] {
	panic("not implemented")
}
