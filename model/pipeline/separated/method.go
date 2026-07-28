// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package separated

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/result"
	"github.com/andreychh/tgen/model/typeform"
)

// Method is the record of a documentation method whose return is split into
// what the method signals: its reference, name, description, and result.
type Method struct {
	Ref         model.Reference
	Name        model.Name
	Description prose.Passage
	Result      result.Result
}

// MethodMapping maps a flattened method into a separated one by deciding
// whether its return type carries data or only reports success.
type MethodMapping struct{}

// NewMethodMapping constructs a MethodMapping.
func NewMethodMapping() MethodMapping {
	return MethodMapping{}
}

// Apply implements [pipeline.Mapping]. It never fails.
func (m MethodMapping) Apply(method flattened.Method) (Method, error) {
	return Method{
		Ref:         method.Ref,
		Name:        method.Name,
		Description: method.Description,
		Result:      m.classify(method.Type),
	}, nil
}

// classify returns [result.Confirmation] when the method returns nothing but
// the literal True, and [result.Value] otherwise. An array of True and a named
// type that happens to hold True are values: only the bare keyword signals
// success and carries no data.
func (m MethodMapping) classify(typ typeform.Type) result.Result {
	atom, ok := typ.Atom().(typeform.Primitive)
	if !ok || atom.Kind() != primitive.True || typ.Dimensionality() != 0 {
		return result.NewValue(typ)
	}
	return result.NewConfirmation()
}
