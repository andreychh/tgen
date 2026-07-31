// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

// Payload represents the request slot of a generated method: the concrete type
// its payload method returns, and the template that renders that method's body.
// The variants are exclusive — a method assembles its request exactly one way.
//
//sumtype:decl
type Payload interface {
	// Type returns the payload type the payload method hands back.
	Type() string
	// Template returns the name of the template rendering the payload body.
	Template() string

	isPayload()
}

// Empty represents the request of a method that takes no parameters, carrying
// neither a body nor a content type.
type Empty struct{}

// NewEmpty creates an Empty.
func NewEmpty() Empty {
	return Empty{}
}

// Type implements [Payload].
func (Empty) Type() string {
	return "emptyPayload"
}

// Template implements [Payload].
func (Empty) Template() string {
	return "payload_empty"
}

func (Empty) isPayload() {}

// JSON represents the request of a method that reaches no file, marshalled
// whole from the method struct.
type JSON struct{}

// NewJSON creates a JSON.
func NewJSON() JSON {
	return JSON{}
}

// Type implements [Payload].
func (JSON) Type() string {
	return "jsonPayload"
}

// Template implements [Payload].
func (JSON) Template() string {
	return "payload_json"
}

func (JSON) isPayload() {}

// Form represents the request of a method reaching a file. A file cannot travel
// inside JSON, so every file the parameters reach is collected into a part of
// its own and the body is rewritten to point at those parts.
type Form struct {
	files []Placed
}

// NewForm creates a Form from the parameters reaching a file.
func NewForm(files []Placed) Form {
	return Form{files: files}
}

// Type implements [Payload].
func (Form) Type() string {
	return "formPayload"
}

// Template implements [Payload].
func (Form) Template() string {
	return "payload_form"
}

// Files returns the parameters that hand a file over, in the order the
// documentation listed them.
func (f Form) Files() []Placed {
	return f.files
}

func (Form) isPayload() {}
