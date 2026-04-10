// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model"

type Annotation struct {
	typ      model.Type
	optional model.Optionality
}

func NewAnnotation(t model.Type, o model.Optionality) Annotation {
	return Annotation{typ: t, optional: o}
}

func (a Annotation) AsString() (string, error) {
	typ, err := NewType(a.typ).AsString()
	if err != nil {
		return "", err
	}
	opt, err := a.optional.AsBool()
	if err != nil {
		return "", err
	}
	if !opt {
		return typ, nil
	}
	return typ + " | None = None", nil
}
