// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"github.com/iancoleman/strcase"

	"github.com/andreychh/tgen/model"
	ir "github.com/andreychh/tgen/model/ir/v2"
)

// filed represents what every file field holds whatever hands its file over:
// the declaration it shadows, and the local the handed-over value waits in.
type filed struct {
	inner ir.FileField
}

// Name returns the Go name of the declaration, which the shadow declares too.
func (f filed) Name() string {
	return NewNameFromKey(f.inner.Key).Value()
}

// Key returns the encoded key of the declaration.
func (f filed) Key() string {
	return string(f.inner.Key)
}

// Local returns the name of the local the handed-over value is held in until
// the shadow is assembled.
func (f filed) Local() string {
	return strcase.ToLowerCamel(string(f.inner.Key))
}

// Tag returns the struct tag of the shadow, omitting an optional declaration
// the caller left unset.
func (f filed) Tag() string {
	return NewTag(f.inner.Key, f.inner.Optionality).Value()
}

// Array reports whether the declaration holds a sequence, each element of which
// hands its own file over.
func (f filed) Array() bool {
	return f.inner.Type.Dimensionality() > 0
}

// Optional reports whether the caller may leave the declaration unset, which
// leaves nothing to hand over.
func (f filed) Optional() bool {
	return bool(f.inner.Optionality)
}

// file reports whether the declaration is the file itself rather than something
// holding one inside.
func (f filed) file() bool {
	return f.inner.Kind == model.FileKindFile
}

// carried returns the Go type a declaration holding a file inside is shadowed
// by. It rewrites itself into JSON wherever it sits, so where it sits does not
// change the answer.
func (f filed) carried() string {
	if f.Array() {
		return "[]json.RawMessage"
	}
	return "json.RawMessage"
}

// carrying returns the name of the template handing over the file a declaration
// holds inside, which for the same reason does not depend on where it sits.
func (f filed) carrying() string {
	if f.Array() {
		return "file_resolve_array"
	}
	return "file_resolve"
}

// Placed represents a parameter of a method that reaches a file. A parameter
// owns a key of the request, so its file travels in a part named by that key
// and leaves behind the reference the sink wrote for it.
type Placed struct {
	filed
}

// NewPlaced creates a Placed from the record of a parameter reaching a file.
func NewPlaced(f ir.FileField) Placed {
	return Placed{filed: filed{inner: f}}
}

// Receiver returns the name the enclosing method's receiver is written as.
func (Placed) Receiver() string {
	return "m"
}

// Fail returns the statement propagating err out of the enclosing method.
func (Placed) Fail() string {
	return "return formPayload{}, err"
}

// Shadow returns the Go type standing in for the declared one. Placing a file
// yields a reference only when the caller set one, so it is always a pointer.
func (p Placed) Shadow() string {
	if p.file() {
		return "*string"
	}
	return p.carried()
}

// Tag returns the struct tag of the shadow. A placed file carries nothing in
// the body when the caller uploads it, since the bytes travel in a part named
// by the same key — so the shadow is omitted whenever it is empty, however
// required the parameter is. Writing it out would put a second value under
// that key, which the receiver reads as the file identifier the upload was
// sent instead of.
func (p Placed) Tag() string {
	if p.file() {
		return NewOptionalTag(p.inner.Key).Value()
	}
	return p.filed.Tag()
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

// Receiver returns the name the enclosing method's receiver is written as.
func (Attached) Receiver() string {
	return "o"
}

// Fail returns the statement propagating err out of the enclosing method.
func (Attached) Fail() string {
	return "return nil, err"
}

// Shadow returns the Go type standing in for the declared one. Attaching a file
// always yields a reference, so only an optional field needs a pointer to tell
// an unset one from a reference to nothing.
func (a Attached) Shadow() string {
	if !a.file() {
		return a.carried()
	}
	if a.Optional() {
		return "*string"
	}
	return "string"
}

// Template returns the name of the template handing the file over.
func (a Attached) Template() string {
	if a.file() {
		return "file_attach"
	}
	return a.carrying()
}
