// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"
	"iter"
	"strings"

	"github.com/andreychh/tgen/model/explicit"
)

type Variants struct {
	inner iter.Seq[explicit.DiscriminatedObject]
}

func NewVariants(s iter.Seq[explicit.DiscriminatedObject]) Variants {
	return Variants{inner: s}
}

func (v Variants) Value() (string, error) {
	var names []string
	for o := range v.inner {
		objName, err := o.Name()
		if err != nil {
			return "", fmt.Errorf("variant name: %w", err)
		}
		name, err := NewClassName(objName).Value()
		if err != nil {
			return "", fmt.Errorf("variant name: %w", err)
		}
		names = append(names, name)
	}
	return strings.Join(names, " | "), nil
}
