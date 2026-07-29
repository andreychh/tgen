// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"github.com/andreychh/tgen/model"
	prosetree "github.com/andreychh/tgen/model/prose"
)

// Definition is the decoded record of anything the page names, whatever its
// kind: its reference, name, kind, position, and description. What a
// definition owns — the fields of an object, the parameters of a method, the
// variants of a union — forms tables of its own.
type Definition struct {
	Ref         model.Reference
	Name        model.Name
	Kind        model.DefinitionKind
	Position    model.Position
	Description prosetree.Passage
}

// KindFilter is a [pipeline.Filter] that narrows a definitions table to the
// definitions of one kind.
type KindFilter struct {
	kind model.DefinitionKind
}

// NewKindFilter constructs a KindFilter keeping the definitions of kind.
func NewKindFilter(kind model.DefinitionKind) KindFilter {
	return KindFilter{kind: kind}
}

// Apply implements [pipeline.Filter]. It reports whether definition is of the
// kind the filter keeps.
func (f KindFilter) Apply(_ model.Reference, definition Definition) bool {
	return definition.Kind == f.kind
}
