// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package flattened

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/resolved"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeform"
)

// Method is the record of a documentation method whose return type is reduced
// to a flat one: its reference, name, description, and the type it returns.
type Method struct {
	Ref         model.Reference
	Name        model.Name
	Description prose.Passage
	Type        typeform.Type
}

// MethodMapping maps a resolved method into a flattened one by reducing its
// return type to a flat type.
type MethodMapping struct{}

// NewMethodMapping constructs a MethodMapping.
func NewMethodMapping() MethodMapping {
	return MethodMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the method's return type
// holds a union.
func (m MethodMapping) Apply(method resolved.Method) (Method, error) {
	typ, err := NewNarrowing(method.Type).Value()
	if err != nil {
		return Method{}, fmt.Errorf("narrowing return type: %w", err)
	}
	return Method{
		Ref:         method.Ref,
		Name:        method.Name,
		Description: method.Description,
		Type:        typ,
	}, nil
}
