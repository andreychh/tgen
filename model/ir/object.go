// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// Object represents a Telegram Bot API object definition narrowed for code
// generation.
type Object struct {
	inner spec.Object
}

// NewObject constructs an Object from a parsed object.
func NewObject(o spec.Object) Object {
	return Object{inner: o}
}

func (o Object) Reference() (model.Reference, error) {
	return o.inner.Reference()
}

func (o Object) Name() (model.Name, error) {
	return o.inner.Name()
}

func (o Object) Description() model.Description {
	return o.inner.Description()
}

func (o Object) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(o.inner.Fields(), NewField)
}

func (o Object) HasInputFile() (bool, error) {
	for f := range o.Fields() {
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
