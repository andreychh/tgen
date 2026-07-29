// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package flattened

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/resolved"
	"github.com/andreychh/tgen/model/typeform"
)

// Method is what a method has of its own, with the type it returns reduced to a
// flat one.
type Method struct {
	Ref  model.Reference
	Type typeform.Type
}

// MethodMapping maps a resolved method into a flattened one by reducing the
// type it returns to a flat one.
type MethodMapping struct{}

// NewMethodMapping constructs a MethodMapping.
func NewMethodMapping() MethodMapping {
	return MethodMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the return type holds a
// union.
func (m MethodMapping) Apply(method resolved.Method) (Method, error) {
	typ, err := NewNarrowing(method.Type).Value()
	if err != nil {
		return Method{}, fmt.Errorf("narrowing return type: %w", err)
	}
	return Method{
		Ref:  method.Ref,
		Type: typ,
	}, nil
}
