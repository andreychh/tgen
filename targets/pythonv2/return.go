// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Return represents the response slot of a generated method: what its call
// annotates as the answer, and the template rendering the body reading it.
//
// The Go target parts four variants here, two of them for a union, because
// encoding/json cannot decode into an interface and a method answering with one
// has to hand the payload to a decoder written by hand. Pydantic tells a union
// apart on its own, so what is left is the one distinction the documentation
// itself draws: whether a method answers with a value or only reports that it
// worked.
//
//sumtype:decl
type Return interface {
	// Signature returns what the call annotates as its answer.
	Signature() string
	// Template returns the name of the template rendering the response body.
	Template() string

	isReturn()
}

// Command represents a return that only signals success. The documentation
// writes it as True, which carries nothing a caller could read, so the call
// answers with None and a failure arrives as the exception the connection
// raises.
type Command struct{}

// NewCommand creates a Command.
func NewCommand() Command {
	return Command{}
}

// Signature implements [Return].
func (Command) Signature() string {
	return none
}

// Template implements [Return].
func (Command) Template() string {
	return "return_command"
}

func (Command) isReturn() {}

// Value represents a return carrying something the caller reads.
type Value struct {
	annotation Annotation
}

// NewValue creates a Value from a response type.
func NewValue(annotation Annotation) Value {
	return Value{annotation: annotation}
}

// Signature implements [Return].
func (v Value) Signature() string {
	return v.annotation.Value()
}

// Template implements [Return].
func (Value) Template() string {
	return "return_value"
}

// Adapter returns the adapter validating the response into the type the call
// answers with.
func (v Value) Adapter() string {
	return "TypeAdapter(" + v.annotation.Value() + ")"
}

func (Value) isReturn() {}

// Result represents the Python view of what a method returns.
type Result struct {
	inner ir.Result
}

// NewResult creates a Result from a resolved method result.
func NewResult(r ir.Result) Result {
	return Result{inner: r}
}

// Return returns the [Return] variant rendering the result. It panics on a
// result of a kind the target renders nothing for, the documentation drawing one
// distinction here and this reading both sides of it.
func (r Result) Return() Return {
	switch inner := r.inner.(type) {
	case ir.Confirmation:
		return NewCommand()
	case ir.Value:
		return NewValue(NewRequiredAnnotation(inner.Type()))
	}
	panic(fmt.Sprintf("pythonv2: unknown result %T", r.inner))
}
