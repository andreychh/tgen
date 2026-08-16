// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

// Payload represents the request slot of a generated method: the template
// rendering the body of its payload method. The variants are exclusive — a
// method assembles its request exactly one way.
//
//sumtype:decl
type Payload interface {
	// Template returns the name of the template rendering the payload body.
	Template() string

	isPayload()
}

// Empty represents the request of a method that takes no parameter, carrying
// neither a body nor a content type.
type Empty struct{}

// NewEmpty creates an Empty.
func NewEmpty() Empty {
	return Empty{}
}

// Template implements [Payload].
func (Empty) Template() string {
	return "payload_empty"
}

func (Empty) isPayload() {}

// JSON represents the request of a method that reaches no file, dumped whole
// from the model holding its parameters.
type JSON struct{}

// NewJSON creates a JSON.
func NewJSON() JSON {
	return JSON{}
}

// Template implements [Payload].
func (JSON) Template() string {
	return "payload_json"
}

func (JSON) isPayload() {}
