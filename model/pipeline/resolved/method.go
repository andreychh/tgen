// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package resolved

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/resolved/types"
	"github.com/andreychh/tgen/model/result"
	typetree "github.com/andreychh/tgen/model/types/v2"
)

// Method is the decoded record of a documentation method after its return type
// is resolved: its reference, name, and what it returns.
type Method struct {
	Ref    model.Reference
	Name   model.Name
	Result result.Result
}

// MethodMapping maps a parsed method into a resolved method by decoding its
// return type from description prose into a result.
type MethodMapping struct{}

// NewMethodMapping constructs a MethodMapping.
func NewMethodMapping() MethodMapping {
	return MethodMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the method's description
// prose does not contain a recognizable return clause.
func (m MethodMapping) Apply(method parsed.Method) (Method, error) {
	expr, err := types.NewReturnType(method.Description).Value()
	if err != nil {
		return Method{}, fmt.Errorf("decoding return type: %w", err)
	}
	return Method{
		Ref:    method.Ref,
		Name:   method.Name,
		Result: m.classify(expr),
	}, nil
}

// classify returns [result.Confirmation] when expr is the True primitive, and
// [result.Value] otherwise.
func (m MethodMapping) classify(expr typetree.Expression) result.Result {
	if expr.Equals(typetree.NewPrimitive(typetree.True)) {
		return result.NewConfirmation()
	}
	return result.NewValue(expr)
}
