// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"github.com/andreychh/tgen/model"
	ir "github.com/andreychh/tgen/model/ir/v2"
)

// filed represents what every field reaching a file holds whatever hands its
// file over: the attribute the model is read off by, and the key the value
// travels under. No file the documented API names spells the two differently,
// and they are held apart regardless, since the dump is told what to leave out
// by attribute and told back what to write by key — a file field whose key
// Python reserves would part them with nothing here to change.
type filed struct {
	inner ir.FileField
}

// Name returns the Python attribute the field declares, which is also the name
// the dump is told to leave the field out by.
func (f filed) Name() string {
	return NewFieldName(f.inner.Key).Value()
}

// Key returns the key the field travels under.
func (f filed) Key() string {
	return string(f.inner.Key)
}

// Optional reports whether the caller may leave the field unset, which leaves
// nothing to hand over.
func (f filed) Optional() bool {
	return bool(f.inner.Optionality)
}

// file reports whether the field is the file itself rather than something
// holding one inside.
func (f filed) file() bool {
	return f.inner.Kind == model.FileKindFile
}

// carrying returns the name of the template handing over the file a field holds
// inside. It does not turn on where the field sits: such a field rewrites itself
// into the JSON its owner carries wherever it stands, so a parameter and a field
// hand it over the same way.
func (f filed) carrying() string {
	if f.inner.Type.Dimensionality() > 0 {
		return "file_resolve_array"
	}
	return "file_resolve"
}

// Placed represents a parameter of a method that reaches a file. A parameter
// owns a key of the request, so its file travels in a part named by that key
// and leaves that key of the body empty.
type Placed struct {
	filed
}

// NewPlaced creates a Placed from the record of a parameter reaching a file.
func NewPlaced(f ir.FileField) Placed {
	return Placed{filed: filed{inner: f}}
}

// Template returns the name of the template handing the file over.
func (p Placed) Template() string {
	if p.file() {
		return "file_place"
	}
	return p.carrying()
}

// Attached represents a field of an object that reaches a file. An object sits
// inside a value of the request and owns no key of it, so its file travels in a
// part under a generated key and leaves behind a reference to that key.
type Attached struct {
	filed
}

// NewAttached creates an Attached from the record of a field reaching a file.
func NewAttached(f ir.FileField) Attached {
	return Attached{filed: filed{inner: f}}
}

// Template returns the name of the template handing the file over.
func (a Attached) Template() string {
	if a.file() {
		return "file_attach"
	}
	return a.carrying()
}
