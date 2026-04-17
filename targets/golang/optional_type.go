// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"
	"strings"

	"github.com/andreychh/tgen/model"
)

// OptionalType decorates Type, adding a * prefix for optional non-array non-union types.
type OptionalType struct {
	inner Type
	opt   model.Optionality
}

// NewOptionalType creates an OptionalType from a model.Type and its optionality.
func NewOptionalType(t Type, o model.Optionality) OptionalType {
	return OptionalType{inner: t, opt: o}
}

func (t OptionalType) IsPrimitive() (bool, error) {
	return t.inner.IsPrimitive()
}

func (t OptionalType) IsUnion() (bool, error) {
	return t.inner.IsUnion()
}

func (t OptionalType) Depth() (int, error) {
	return t.inner.Depth()
}

func (t OptionalType) Name() (string, error) {
	return t.inner.Name()
}

func (t OptionalType) Part() (string, error) {
	part, err := t.inner.Part()
	if err != nil {
		return "", err
	}
	opt, err := t.opt.AsBool()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if !opt {
		return part, nil
	}
	depth, err := t.inner.Depth()
	if err != nil {
		return "", err
	}
	if depth > 0 {
		return part, nil
	}
	isUnion, err := t.inner.IsUnion()
	if err != nil {
		return "", err
	}
	if isUnion {
		return part, nil
	}
	isPrimitive, err := t.inner.IsPrimitive()
	if err != nil {
		return "", err
	}
	if isPrimitive {
		return strings.Replace(part, "%s", "*%s", 1), nil
	}
	return part, nil
}

func (t OptionalType) Zero() (string, error) {
	opt, err := t.opt.AsBool()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if opt {
		return zeroNil, nil
	}
	return t.inner.Zero()
}

func (t OptionalType) AsString() (string, error) {
	str, err := t.inner.AsString()
	if err != nil {
		return "", err
	}
	depth, err := t.inner.Depth()
	if err != nil {
		return "", err
	}
	isUnion, err := t.inner.IsUnion()
	if err != nil {
		return "", err
	}
	opt, err := t.opt.AsBool()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if !opt || depth > 0 || isUnion {
		return str, nil
	}
	return "*" + str, nil
}
