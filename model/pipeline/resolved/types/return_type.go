// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"errors"

	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/types/v2"
)

// ReturnType is a method's description ready to be decoded into the type it
// returns.
type ReturnType struct {
	description prose.Passage
}

// NewReturnType constructs a ReturnType over a method's description.
func NewReturnType(description prose.Passage) ReturnType {
	return ReturnType{description: description}
}

// Value returns the type expression decoded from the return clause in the first
// paragraph of the description. It fails when the description has no first
// paragraph or no rule recognizes the return clause.
func (r ReturnType) Value() (types.Expression, error) {
	blocks := r.description.Blocks()
	if len(blocks) == 0 {
		return nil, errors.New("method description has no blocks")
	}
	paragraph, ok := blocks[0].(prose.Paragraph)
	if !ok {
		return nil, errors.New("method description does not open with a paragraph")
	}
	expr, found := NewSearch(rules()).Find(paragraph.Inlines())
	if !found {
		return nil, errors.New("no rule matched the return clause")
	}
	return expr, nil
}

// rules assembles the return-clause production in priority order: the specific
// array and union forms before the plain named and primitive forms, so a clause
// like "Array of Update" is not first claimed by a plain named or primitive
// rule.
func rules() ProductionRule {
	return NewProductionRule(
		NewUnionRule(),
		NewArrayRule(),
		NewReturnsRule(),
		NewReturnedRule(),
	)
}
