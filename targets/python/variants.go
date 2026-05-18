// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"
	"iter"
	"strings"

	"github.com/andreychh/tgen/model/ir"
)

type Variants struct {
	inner iter.Seq[ir.DiscriminatedObject]
}

func NewVariants(s iter.Seq[ir.DiscriminatedObject]) Variants {
	return Variants{inner: s}
}

func (v Variants) Value() (string, error) {
	var names []string
	for o := range v.inner {
		name, err := o.Name()
		if err != nil {
			return "", fmt.Errorf("variant name: %w", err)
		}
		names = append(names, NewClassName(name).Value())
	}
	return strings.Join(names, " | "), nil
}
