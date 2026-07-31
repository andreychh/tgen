// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/separated"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/result"
)

// Method is the record of a documented method, joined with the parameters it
// takes and the result it returns. Name is the name the endpoint is called by,
// which is also the name a target derives its identifier from. Files narrows
// Params to those reaching a file, and is empty for a method sending none.
type Method struct {
	Ref         model.Reference
	Name        model.Name
	Description prose.Passage
	Params      []Field
	Files       []FileField
	Result      Result
}

func (Method) isDefinition() {}

// Returned is the view of what one method returns, with the type it carries
// bound to the definition that type names.
type Returned struct {
	db    separated.Specification
	inner result.Result
}

// NewReturned constructs a Returned over the result a method was separated
// into.
func NewReturned(db separated.Specification, inner result.Result) Returned {
	return Returned{db: db, inner: inner}
}

// Value returns the result with its carried type bound. It fails when that type
// names no definition.
func (r Returned) Value() (Result, error) {
	switch inner := r.inner.(type) {
	case result.Confirmation:
		return NewConfirmation(), nil
	case result.Value:
		typ, err := NewResolution(r.db, inner.Type()).Value()
		if err != nil {
			return nil, fmt.Errorf("binding the returned type: %w", err)
		}
		return NewValue(typ), nil
	default:
		return nil, fmt.Errorf("unknown result %T", inner)
	}
}
