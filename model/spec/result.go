// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package spec

import "github.com/andreychh/tgen/model/types"

// Result represents what a Telegram Bot API method returns: either a Command
// that only signals success or a Value that carries a typed return expression.
//
//sumtype:decl
type Result interface {
	isResult()
}

// Command represents a method result that only signals success.
type Command struct{}

// NewCommand constructs a Command.
func NewCommand() Command {
	return Command{}
}

func (Command) isResult() {}

// Value represents a method result that carries a typed return expression.
type Value struct {
	expr types.Expression
}

// NewValue constructs a Value from a return type expression.
func NewValue(expr types.Expression) Value {
	return Value{expr: expr}
}

// Type returns the return type expression carried by the Value.
func (v Value) Type() types.Expression {
	return v.expr
}

func (Value) isResult() {}
