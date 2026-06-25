// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
)

// typeNamePattern matches a type name: a single PascalCase token (object or
// union).
var typeNamePattern = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)

// TypeName is the title cell of a type's <h4> header, naming an object or union
// in PascalCase.
type TypeName struct {
	h4 *goquery.Selection
}

// NewTypeName constructs a TypeName over a type's <h4> header.
func NewTypeName(h4 *goquery.Selection) TypeName {
	return TypeName{h4: h4}
}

// Value returns the type's name. It fails when the title is not a single
// PascalCase token.
func (n TypeName) Value() (model.Name, error) {
	name := n.h4.Text()
	if !typeNamePattern.MatchString(name) {
		return "", fmt.Errorf("type name %q is not a PascalCase token", name)
	}
	return model.Name(name), nil
}

// methodNamePattern matches a method name: a single camelCase token.
var methodNamePattern = regexp.MustCompile(`^[a-z][A-Za-z]*$`)

// MethodName is the title cell of a method's <h4> header, naming the method in
// camelCase.
type MethodName struct {
	h4 *goquery.Selection
}

// NewMethodName constructs a MethodName over a method's <h4> header.
func NewMethodName(h4 *goquery.Selection) MethodName {
	return MethodName{h4: h4}
}

// Value returns the method's name. It fails when the title is not a single
// camelCase token.
func (n MethodName) Value() (model.Name, error) {
	name := n.h4.Text()
	if !methodNamePattern.MatchString(name) {
		return "", fmt.Errorf("method name %q is not a camelCase token", name)
	}
	return model.Name(name), nil
}
