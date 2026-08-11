// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// Method represents the Go declaration of a documented method: the struct
// holding its parameters, the payload that struct is sent as, and the response
// its Call decodes.
type Method struct {
	inner ir.Method
}

// NewMethod creates a Method from the record of a method.
func NewMethod(m ir.Method) Method {
	return Method{inner: m}
}

// Doc returns the doc comment of the declaration, closing with a link back to
// the section the method was read from.
func (m Method) Doc() string {
	return NewDefinitionDoc(m.inner.Ref, m.inner.Description, m.inner.Introduced).Value()
}

// Ref implements [Declaration].
func (m Method) Ref() string {
	return string(m.inner.Ref)
}

// Template implements [Declaration].
func (m Method) Template() string {
	return "method"
}

// Name returns the Go name the method declares: the documented name suffixed
// with Method, which keeps the struct holding a request apart from the object
// the request answers with.
func (m Method) Name() string {
	return NewName(m.inner.Name).Value() + "Method"
}

// Wire returns the name the endpoint is called by.
func (m Method) Wire() string {
	return string(m.inner.Name)
}

// Fields returns the parameters the method takes, in the order the
// documentation listed them.
func (m Method) Fields() []Field {
	return slices.NewMapped(m.inner.Params, NewField)
}

// Return returns the response slot of the method.
func (m Method) Return() Return {
	return NewResult(m.inner.Result).Return()
}

// Payload returns the request slot of the method: nothing to send when it takes
// no parameter, a multipart body when a parameter reaches a file, and plain
// JSON otherwise.
func (m Method) Payload() Payload {
	if len(m.inner.Params) == 0 {
		return NewEmpty()
	}
	if len(m.inner.Files) > 0 {
		return NewForm(slices.NewMapped(m.inner.Files, NewPlaced))
	}
	return NewJSON()
}
