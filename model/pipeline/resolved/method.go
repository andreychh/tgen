// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package resolved

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/resolved/types"
	"github.com/andreychh/tgen/model/typeexpr"
)

// Method is what a method has of its own: the type it returns. Everything it
// shares with any other definition — its name, position, and description —
// stays in the definitions table.
type Method struct {
	Ref  model.Reference
	Type typeexpr.Expression
}

// MethodMapping maps a method definition into what the method has of its own by
// decoding the return type out of its description prose.
type MethodMapping struct{}

// NewMethodMapping constructs a MethodMapping.
func NewMethodMapping() MethodMapping {
	return MethodMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the definition's
// description prose does not contain a recognizable return clause.
func (m MethodMapping) Apply(definition parsed.Definition) (Method, error) {
	expr, err := types.NewReturnType(definition.Description).Value()
	if err != nil {
		return Method{}, fmt.Errorf("decoding return type: %w", err)
	}
	return Method{
		Ref:  definition.Ref,
		Type: expr,
	}, nil
}
