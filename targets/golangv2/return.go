// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Return represents the response slot of a generated method: the result list of
// its Call, the statement that propagates a failure out of Call, and the
// template that renders the response body.
//
//sumtype:decl
type Return interface {
	// Signature returns the result list of Call, parenthesised when the method
	// carries a value.
	Signature() string
	// Fail returns the statement propagating err out of Call.
	Fail() string
	// Template returns the name of the template rendering the response body.
	Template() string

	isReturn()
}

// Command represents a return that only signals success.
type Command struct{}

// NewCommand creates a Command.
func NewCommand() Command {
	return Command{}
}

// Signature implements [Return].
func (Command) Signature() string {
	return "error"
}

// Fail implements [Return].
func (Command) Fail() string {
	return "return err"
}

// Template implements [Return].
func (Command) Template() string {
	return "return_command"
}

func (Command) isReturn() {}

// typed represents the response type shared by every [Return] that carries one.
type typed struct {
	typ Type
}

// Signature implements [Return].
func (t typed) Signature() string {
	return "(" + t.typ.Value() + ", error)"
}

// Fail implements [Return].
func (t typed) Fail() string {
	return "return " + t.typ.Zero() + ", err"
}

// Value returns the rendered response type.
func (t typed) Value() string {
	return t.typ.Value()
}

// Name returns the Go name of the response type, without the slices enclosing
// it.
func (t typed) Name() string {
	return t.typ.Name()
}

// Plain represents a return decoded straight into its Go type.
type Plain struct {
	typed
}

// NewPlain creates a Plain from a response type.
func NewPlain(typ Type) Plain {
	return Plain{typed: typed{typ: typ}}
}

// Template implements [Return].
func (Plain) Template() string {
	return "return_plain"
}

func (Plain) isReturn() {}

// Dispatched represents a return decoded raw and handed to the unmarshaller of
// a union, since encoding/json cannot decode into an interface on its own.
type Dispatched struct {
	typed
}

// NewDispatched creates a Dispatched from a response type.
func NewDispatched(typ Type) Dispatched {
	return Dispatched{typed: typed{typ: typ}}
}

// Template implements [Return].
func (Dispatched) Template() string {
	return "return_union"
}

func (Dispatched) isReturn() {}

// DispatchedArray represents a return handed to the unmarshaller of a union one
// element at a time.
type DispatchedArray struct {
	typed
}

// NewDispatchedArray creates a DispatchedArray from a response type.
func NewDispatchedArray(typ Type) DispatchedArray {
	return DispatchedArray{typed: typed{typ: typ}}
}

// Template implements [Return].
func (DispatchedArray) Template() string {
	return "return_union_array"
}

func (DispatchedArray) isReturn() {}

// Result represents the Go view of what a method returns.
type Result struct {
	inner ir.Result
}

// NewResult creates a Result from a resolved method result.
func NewResult(r ir.Result) Result {
	return Result{inner: r}
}

// Return returns the [Return] variant rendering the result.
func (r Result) Return() Return {
	switch inner := r.inner.(type) {
	case ir.Confirmation:
		return NewCommand()
	case ir.Value:
		return r.carried(NewRequiredType(inner.Type()))
	default:
		panic(fmt.Sprintf("golangv2: unknown result %T", inner))
	}
}

// carried returns the variant rendering a result that carries a type: a union
// is decoded raw and dispatched, one element at a time when it arrives as an
// array, and everything else decodes straight into its Go type.
func (r Result) carried(typ Type) Return {
	if !typ.Union() {
		return NewPlain(typ)
	}
	if typ.Array() {
		return NewDispatchedArray(typ)
	}
	return NewDispatched(typ)
}
