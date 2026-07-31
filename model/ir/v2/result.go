// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model/typebound"
)

// Result represents what a method returns, with the type it carries bound to
// the definition it names. The split between signalling and carrying is decided
// in the pipeline; only the type is joined here. The concrete variants are
// [Confirmation] and [Value].
//
//sumtype:decl
type Result interface {
	isResult()
}

// Confirmation represents a method result that only signals success.
type Confirmation struct{}

// NewConfirmation constructs a Confirmation.
func NewConfirmation() Confirmation {
	return Confirmation{}
}

func (Confirmation) isResult() {}

// Value represents a method result that carries the resolved type it returns.
type Value struct {
	typ typebound.Type
}

// NewValue constructs a Value from the resolved type a method returns.
func NewValue(typ typebound.Type) Value {
	return Value{typ: typ}
}

// Type returns the bound type the Value carries.
func (v Value) Type() typebound.Type {
	return v.typ
}

func (Value) isResult() {}
