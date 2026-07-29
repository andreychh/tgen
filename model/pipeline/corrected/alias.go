// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/typeexpr"
)

// Alias is the type a name tgen introduces stands for. What the alias is
// called and how it is described lives in the definitions table under the
// alias kind, alongside everything else the pipeline names.
type Alias struct {
	Ref  model.Reference
	Type typeexpr.Expression
}
