// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import "github.com/andreychh/tgen/parsing/gq"

type Variant struct {
	selection gq.Selection
}

func NewVariant(li gq.Selection) Variant {
	return Variant{selection: li}
}

func (v Variant) Name() ObjectName {
	return NewObjectName(v.selection.Find("a"))
}
