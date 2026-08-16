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
// rendered as. A kind the target has yet to spell is rendered as a [Stub], so
// the sequence a file walks is the whole page from the first run and a kind
// that lands is a diff of the output rather than the first output there is.
func NewDeclaration(record ir.Definition) Declaration {
	switch record := record.(type) {
	case ir.Object:
		return NewStub(record.Ref, "object")
	case ir.DiscriminatedObject:
		return NewStub(record.Ref, "discriminated object")
	case ir.Union:
		return NewStub(record.Ref, "union")
	case ir.DiscriminatedUnion:
		return NewStub(record.Ref, "discriminated union")
	case ir.Alias:
		return NewStub(record.Ref, "alias")
	case ir.Method:
		return NewStub(record.Ref, "method")
	}
	panic(fmt.Sprintf("pythonv2: unknown definition %T", record))
}
