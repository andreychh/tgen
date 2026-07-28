// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

// Alias is a name tgen introduces for a type expression the documentation
// leaves unnamed. Its description is tgen's own, since the documentation has
// none for a type it never names.
type Alias struct {
	Ref         model.Reference
	Name        model.Name
	Type        typeexpr.Expression
	Description prose.Passage
}
