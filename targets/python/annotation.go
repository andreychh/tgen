// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
)

type Annotation struct {
	typ      types.Expression
	optional model.Optionality
}

func NewAnnotation(t types.Expression, o model.Optionality) Annotation {
	return Annotation{typ: t, optional: o}
}

func (a Annotation) Value() (string, error) {
	typ, err := NewType(a.typ).Value()
	if err != nil {
		return "", err
	}
	if !bool(a.optional) {
		return typ, nil
	}
	return typ + " | None = None", nil
}
