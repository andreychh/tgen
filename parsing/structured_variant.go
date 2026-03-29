// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQStructuredVariant struct {
	li gq.Selection
}

func NewGQStructuredVariant(li gq.Selection) GQStructuredVariant {
	return GQStructuredVariant{li: li}
}

func (v GQStructuredVariant) Reference() Reference {
	// TODO implement me
	panic("implement me")
}

func (v GQStructuredVariant) Name() Name {
	// TODO implement me
	panic("implement me")
}

func (v GQStructuredVariant) Description() Description {
	// TODO implement me
	panic("implement me")
}

func (v GQStructuredVariant) Fields() iter.Seq[Field] {
	// TODO implement me
	panic("implement me")
}
