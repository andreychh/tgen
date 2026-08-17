// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Declaration represents one declaration of the generated package. The file
// walking the sequence knows three things about it: the template rendering its
// shape, the reference a block written by hand claims it by, and the name it
// declares. What that template reads beyond them is the business of the view
// behind it, and no two views need have anything else in common.
//
// The name is here because the package states its own surface, and stating it
// means naming every declaration outside the file that declares them. A template
// reaches the name of a concrete view without being told the interface has one,
// so this could have gone unwritten — and then a kind that stopped answering
// would be a template failing at generation time instead of a package failing to
// compile.
type Declaration interface {
	// Ref returns the reference the declaration is addressed by. A block written
	// by hand claims the declaration by that reference, which no target respells.
	Ref() string
	// Template returns the name of the template rendering the declaration.
	Template() string
	// Name returns the Python name the declaration declares. Every kind declares
	// exactly one, whatever else it holds.
	Name() string
}

// NewDeclaration creates the declaration one record of the pipeline's exit is
// rendered as. It panics on a record of a kind the target renders nothing for,
// every kind the pipeline can hand over being spelled here.
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
