// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

// Result represents what a Telegram Bot API method returns: either a Command
// that only signals success or a Value that carries a resolved return type.
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

// Value represents a method result that carries a resolved return type.
type Value struct {
	typ Type
}

// NewValue constructs a Value from a resolved return type.
func NewValue(typ Type) Value {
	return Value{typ: typ}
}

// Type returns the resolved return type carried by the Value.
func (v Value) Type() Type {
	return v.typ
}

func (Value) isResult() {}
