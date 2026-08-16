// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Field represents the Python declaration of a field an object owns or a
// parameter a method takes.
type Field struct {
	inner ir.Field
}

// NewField creates a Field from the record of a field.
func NewField(f ir.Field) Field {
	return Field{inner: f}
}

// Doc returns the docstring of the declaration, which stands under it rather
// than over it, since a docstring is read from what precedes it.
func (f Field) Doc() string {
	return NewFieldDocstring(f.inner.Description).Value()
}

// Name returns the Python attribute the field declares.
func (f Field) Name() string {
	return NewFieldName(f.inner.Key).Value()
}

// Annotation returns the Python type annotation of the field.
func (f Field) Annotation() string {
	return NewAnnotation(f.inner.Type, f.inner.Optionality).Value()
}

// Assignment returns what the declaration stands equal to: the None an optional
// field defaults to, and the key it is read from where the attribute it
// declares stopped spelling that key. It is empty for a required field whose
// attribute spells its key, which is what leaves the annotation standing alone.
//
// The key is written per field rather than derived for every field at once. One
// key of the documented API — from — is a word Python reserves, so a rule
// covering every model would be a rule the whole API pays for and one field
// uses, and it would misread the first key that ever ends in an underscore of
// its own.
func (f Field) Assignment() string {
	alias := f.alias()
	if alias == "" {
		if f.inner.Optionality {
			return "None"
		}
		return ""
	}
	if f.inner.Optionality {
		return fmt.Sprintf("Field(default=None, alias=%q)", alias)
	}
	return fmt.Sprintf("Field(alias=%q)", alias)
}

// Aliased reports whether the field is read from a key its attribute stopped
// spelling, which is what obliges the class declaring it to admit both names.
func (f Field) Aliased() bool {
	return f.alias() != ""
}

// alias returns the key the field is read from, empty where the attribute it
// declares already spells that key.
func (f Field) alias() string {
	if f.Name() == string(f.inner.Key) {
		return ""
	}
	return string(f.inner.Key)
}
