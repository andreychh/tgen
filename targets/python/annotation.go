// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/ir"
)

type Annotation struct {
	typ ir.Type
	opt model.Optionality
}

func NewAnnotation(t ir.Type, o model.Optionality) Annotation {
	return Annotation{typ: t, opt: o}
}

func (a Annotation) Value() (string, error) {
	typ, err := NewType(a.typ).Value()
	if err != nil {
		return "", err
	}
	if !a.opt {
		return typ, nil
	}
	return typ + " | None = None", nil
}
