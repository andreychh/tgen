// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Declaration represents one declaration of the generated package. The file
// walking the sequence knows only two things about it: the template rendering
// its shape and the reference a block written by hand claims it by. What that
// template reads is the business of the view behind it, and no two views need
// have anything else in common.
type Declaration interface {
	// Ref returns the reference the declaration is addressed by. A block written
	// by hand claims the declaration by that reference, which no target respells.
	Ref() string
	// Template returns the name of the template rendering the declaration.
	Template() string
}

// NewDeclaration creates the declaration one record of the pipeline's exit is
// rendered as.
func NewDeclaration(record ir.Definition) Declaration {
	switch record := record.(type) {
	case ir.Object:
		return NewObject(record)
	case ir.DiscriminatedObject:
		return NewDiscriminatedObject(record)
	case ir.Union:
		return NewUnion(record)
	case ir.DiscriminatedUnion:
		return NewDiscriminatedUnion(record)
	case ir.Alias:
		return NewAlias(record)
	case ir.Method:
		return NewMethod(record)
	}
	panic(fmt.Sprintf("pythonv2: unknown definition %T", record))
}
