// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package modified

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	typetree "github.com/andreychh/tgen/model/types/v2"
)

// Alias is a name tgen introduces for a type expression the documentation
// leaves unnamed. Its description is tgen's own, since the documentation has
// none for a type it never names.
type Alias struct {
	Ref         model.Reference
	Name        model.Name
	Type        typetree.Expression
	Description prose.Passage
}
