// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
)

// DiscriminatedObject represents a variant of a discriminated union narrowed
// for code generation.
type DiscriminatedObject struct {
	inner spec.DiscriminatedObject
}

// NewDiscriminatedObject constructs a DiscriminatedObject from a parsed variant.
func NewDiscriminatedObject(d spec.DiscriminatedObject) DiscriminatedObject {
	return DiscriminatedObject{inner: d}
}

func (d DiscriminatedObject) Reference() (model.Reference, error) {
	return d.inner.Reference()
}

func (d DiscriminatedObject) Name() (model.Name, error) {
	return d.inner.Name()
}

func (d DiscriminatedObject) Description() model.Description {
	return d.inner.Description()
}

func (d DiscriminatedObject) Fields() Fields {
	return NewFields(d.inner.Fields())
}

func (d DiscriminatedObject) HasInputFile() (bool, error) {
	for f := range d.Fields().Free() {
		ok, err := f.IsInputFile()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
