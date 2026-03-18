// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import "github.com/andreychh/tgen/parsing/gq"

type GQVariant struct {
	selection gq.Selection
}

func NewVariant(li gq.Selection) GQVariant {
	return GQVariant{selection: li}
}

func (v GQVariant) Name() ObjectName {
	return NewGQObjectName(v.selection.Find("a"))
}
