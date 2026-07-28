// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package result models what a Telegram Bot API method returns as one of two
// variants: a [Confirmation] that only signals success, or a [Value] that
// carries the type it returns.
package result

import (
	"github.com/andreychh/tgen/model/typeform"
)

// Result represents what a Telegram Bot API method returns. The concrete
// variants are [Confirmation] and [Value].
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

// Value represents a method result that carries the type it returns.
type Value struct {
	typ typeform.Type
}

// NewValue constructs a Value from a returned type.
func NewValue(typ typeform.Type) Value {
	return Value{typ: typ}
}

// Type returns the type carried by the Value.
func (v Value) Type() typeform.Type {
	return v.typ
}

func (Value) isResult() {}
